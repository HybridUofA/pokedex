package mapper

import (

    "io"
    "net/http"
    "encoding/json"
)

type LocationData struct {
    Next     *string    `json:"next"`
    Previous *string       `json:"previous"`
    Results  []Results `json:"results"`
}
type Results struct {
    Name string `json:"name"`
}

func MapLocations(url string) (LocationData, error) {

    res, err := http.Get(url)
    if err != nil {
        return LocationData{}, err
    }
    defer res.Body.Close()
    body, err := io.ReadAll(res.Body)
    if err != nil {
        return LocationData{}, err
    }

    var locData LocationData
    if err := json.Unmarshal(body, &locData); err != nil {
        return LocationData{}, err
    }

    return locData, nil
}
