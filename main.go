package main

import (
    "fmt"
    "bufio"
    "os"
    "time"
    "encoding/json"
    "github.com/HybridUofA/pokedex/internal/cleaninput"
    "github.com/HybridUofA/pokedex/internal/mapper"
    "github.com/HybridUofA/pokedex/internal/config"
    "github.com/HybridUofA/pokedex/internal/pokecache"
    "github.com/HybridUofA/pokedex/internal/encounter"
    "github.com/HybridUofA/pokedex/internal/pokemon"
)

var pokeCache *pokecache.Cache

type cliCommand struct {
    name        string
    description string
    callback    func(*config.Config, string) error
}

var commands map[string]cliCommand
var pokedex map[string]pokemon.Pokemon

func init() {
    commands = map[string]cliCommand{
        "exit": {
            name:        "exit",
            description: "Exit the Pokedex",
            callback:    commandExit,
        },
        "help": {
            name:        "help",
            description: "Displays a help message",
            callback:    commandHelp,
        },
        "map": {
            name:        "map",
            description: "Displays the next 20 locations on the map",
            callback:    commandMap,
        },
        "mapb": {
            name:        "mapb",
            description: "Displays the previous 20 locations on the map",
            callback:    commandMapBack,
        },
        "explore": {
            name:        "explore",
            description: "Displays the pokemon a player can encounter in the area",
            callback:    commandExplore,
        },
        "catch": {
            name:        "catch",
            description: "Attempts to throw a Pokeball and capture a Pokemon!",
            callback:    commandCatch,
        },
        "inspect": {
            name:        "inspect",
            description: "Inspects the pokedex entry of a specified pokemon",
            callback:    commandInspect,
        },
        "pokedex": {
            name:        "pokedex",
            description: "Displays a list of the pokemon you have caught",
            callback:    commandPokedex,
        },
    }
    pokeCache = pokecache.NewCache(1 * time.Minute)
    pokedex = make(map[string]pokemon.Pokemon)
}


func commandHelp(config *config.Config, arg string) error {
    fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
    for _, command := range commands {
        fmt.Printf("%s: %s\n", command.name, command.description)
    }
    return nil
}

func commandExit(config *config.Config, arg string) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

func commandMap(config *config.Config, arg string) error {

    url := "https://pokeapi.co/api/v2/location-area"
    if config.Next != nil {
        url = *config.Next
    }

    if data, found := pokeCache.Get(url); found {
        var locData mapper.LocationData
        if err := json.Unmarshal(data, &locData); err == nil {
            for _, result := range locData.Results {
                fmt.Println(result.Name)
            }
            config.Next = locData.Next
            config.Previous = locData.Previous
            return nil
        }
    }

    locData, err := mapper.MapLocations(url)
    if err != nil {
        return err
    }
    cacheData, err := json.Marshal(locData)
    if err != nil {
        return err
    }
    pokeCache.Add(url, cacheData)

    for _, result := range locData.Results {
        fmt.Println(result.Name)
    }

    config.Next = locData.Next
    config.Previous = locData.Previous
    return nil
}

func commandMapBack(config *config.Config, arg string) error {

    if config.Previous == nil {
        fmt.Println("you're on the first page")
        return nil
    }

    url := *config.Previous

    if data, found := pokeCache.Get(url); found {
        var locData mapper.LocationData
        if err := json.Unmarshal(data, &locData); err == nil {
            for _, result := range locData.Results {
                fmt.Println(result.Name)
            }
            config.Next = locData.Next
            config.Previous = locData.Previous
        }
        return nil
    }
    locData, err := mapper.MapLocations(url)
    if err != nil {
        return err
    }

    cacheData, err := json.Marshal(locData)
    if err != nil {
        return err
    }
    pokeCache.Add(url, cacheData)

    for _, result := range locData.Results {
        fmt.Println(result.Name)
    }

    config.Next = locData.Next
    config.Previous = locData.Previous
    return nil
}

func commandExplore(cfg *config.Config, arg string) error {
    url := "https://pokeapi.co/api/v2/location-area/" + arg

    fmt.Printf("Exploring %s...\n", arg)
    fmt.Println("Found Pokemon:")

    if data, found := pokeCache.Get(url); found {
        var encData encounter.EncounterData
        if err := json.Unmarshal(data, &encData); err == nil {
            for _, enc := range encData.Encounters {
                fmt.Println(enc.Pokemon.Name)
            }
            return nil
        }
    }

    encData, err := encounter.GetEncounters(url)
    if err != nil {
        return err
    }

    cacheData, err := json.Marshal(encData)
    if err != nil {
        return err
    }
    pokeCache.Add(url, cacheData)

    for _, enc := range encData.Encounters {
        fmt.Println(enc.Pokemon.Name)
    }
    return nil
}

func commandCatch(cfg *config.Config, arg string) error {
    url := "https://pokeapi.co/api/v2/pokemon/" + arg

    fmt.Printf("Throwing a Pokeball at %s...\n", arg)
    pokeData, success, err := pokemon.CatchPokemon(url)
    if err != nil {
        return err
    }
    if success {
        fmt.Printf("%s was caught!\n", arg)
        pokedex[arg] = pokeData
        return nil
    }
    fmt.Printf("%s escaped!\n", arg)
    return nil
}

func commandInspect(cfg *config.Config, arg string) error {
    if p, exists := pokedex[arg]; exists {
        fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\nStats: \n  -hp: %d\n  -attack: %d\n  -defense: %d\n  -special-attack: %d\n  -special-defense: %d\n  -speed: %d\nTypes:\n", p.Name, p.Height, p.Weight, pokemon.GetStatValue(p, "hp"),pokemon.GetStatValue(p, "attack"), pokemon.GetStatValue(p, "defense"), pokemon.GetStatValue(p, "special-attack"), pokemon.GetStatValue(p, "special-defense"), pokemon.GetStatValue(p, "speed"))
        for _, pType := range pokemon.GetTypes(p) {
            fmt.Printf("  -%s\n", pType)
        }
    } else {
    fmt.Println("you have not caught that pokemon")
    }
    return nil
}

func commandPokedex(cgf *config.Config, arg string) error {
    fmt.Println("Your Pokedex:")
    for _, p := range pokedex {
        fmt.Printf(" - %s\n", p.Name)
    }
    return nil
}

func main() {

    config := &config.Config{
        Next:     nil,
        Previous: nil,
    }

    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("Pokedex > ")
        scanner.Scan()
        input := scanner.Text()
        words := cleaninput.CleanInput(input)
        if len(words) > 0 {
            arg := ""
            command := words[0]
            if len(words) > 1 {
                arg = words[1]
            }
            if value, exists := commands[command]; exists {
                err := value.callback(config, arg)
                if err != nil {
                    fmt.Printf("An error has occurred: %v\n", err)
                }
            } else {
                fmt.Println("Unknown command")
            }
        }
    }
}
