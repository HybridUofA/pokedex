package cleaninput

import (
    "strings"
    "unicode"
)

func CleanInput(text string) []string { 
    cleaned := strings.TrimSpace(text)
    cleaned = strings.ToLower(cleaned)

    words := strings.Fields(cleaned)

    var result []string

    for _, word := range words {
        cleanedWord := ""
        for _, r := range word {
            if unicode.IsLetter(r) {
                cleanedWord += string(r)
            }
        }
        result = append(result, cleanedWord)
    }
    return result
}
