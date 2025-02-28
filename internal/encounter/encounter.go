package encounter

import (
    "encoding/json"
    "net/http"
    "io"
)

type EncounterData struct {
    Name        string      `json:"name"`
    Encounters  []Encounter `json:"pokemon_encounters"`
}

type Encounter struct {
    Pokemon struct {
        Name           string `json:"name"`
        URL            string `json:"url"`
        BaseExperience int    `json:"base_experience"`
    } `json:"pokemon"`
}

func GetEncounters(url string) (EncounterData, error) {
    resp, err := http.Get(url)
    if err != nil {
        return EncounterData{}, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return EncounterData{}, err
    }

    var encData EncounterData
    err = json.Unmarshal(body, &encData)
    if err != nil {
        return EncounterData{}, err
    }

    return encData, nil
}
