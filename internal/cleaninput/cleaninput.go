package cleaninput

import (
    "strings"
)

func CleanInput(text string) []string {
    cleaned := strings.TrimSpace(text)
    cleaned = strings.ToLower(cleaned)

    words := strings.Fields(cleaned)

    return words
}
