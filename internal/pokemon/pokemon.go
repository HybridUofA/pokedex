package pokemon

import (
    "net/http"
    "encoding/json"
    "io"
    "time"
    "math/rand"
)

type TypeInfo struct {
    Name string `json:"name"`
}

type PokemonType struct {
    Type TypeInfo `json:"type"`
}

type StatInfo struct {
    Name string `json:"name"`
}

type PokemonStat struct {
    BaseStat int      `json:"base_stat"`
    Stat     StatInfo `json:"stat"`
}

type Pokemon struct {
    Name           string        `json:"name"`
    URL            string        `json:"url"`
    BaseExperience int           `json:"base_experience"`
    Height         int           `json:"height"`
    Weight         int           `json:"weight"`
    Types          []PokemonType `json:"types"`
    Stats          []PokemonStat `json:"stats"`
}

func CatchPokemon(url string) (Pokemon, bool, error) {

    rand.Seed(time.Now().UnixNano())

    res, err := http.Get(url)
    if err != nil {
        return Pokemon{}, false, err
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        return Pokemon{}, false, err
    }

    var pokemon Pokemon
    err = json.Unmarshal(body, &pokemon)
    if err != nil {
        return Pokemon{}, false, err
    }

    catchChance := 100.0 / float64(pokemon.BaseExperience)
    r := rand.Float64()

    catchSuccess := r < catchChance

    return pokemon, catchSuccess, nil
}

func GetTypes(pokemon Pokemon) []string {
    typeNames := []string{}
    for _, pokemonType := range pokemon.Types {
        typeNames = append(typeNames, pokemonType.Type.Name)
    }
    return typeNames
}

func GetStatValue(p Pokemon, statName string) int {
    for _, stat := range p.Stats {
        if stat.Stat.Name == statName {
            return stat.BaseStat
        }
    }
    return 0
}
