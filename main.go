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
)

var pokeCache *pokecache.Cache

type cliCommand struct {
    name        string
    description string
    callback    func(*config.Config) error
}

var commands map[string]cliCommand

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
    }
    pokeCache = pokecache.NewCache(1 * time.Minute)
}


func commandHelp(config *config.Config) error {
    fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
    for _, command := range commands {
        fmt.Printf("%s: %s\n", command.name, command.description)
    }
    return nil
}

func commandExit(config *config.Config) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

func commandMap(config *config.Config) error {

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

func commandMapBack(config *config.Config) error {

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
            command := words[0]
            if value, exists := commands[command]; exists {
                err := value.callback(config)
                if err != nil {
                    fmt.Printf("An error has occurred: %v", err)
                }
            } else {
                fmt.Println("Unknown command")
            }
        }
    }
}
