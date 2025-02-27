package cleaninput

import (
    "testing"
)

func TestCleanInput(t *testing.T) {
    cases := []struct {
        input    string
        expected []string
    }{
        {
            input:   "  hello  world  ",
            expected: []string{"hello", "world"},
        },
        {
            input:   "Charmander Bulbasaur PIKACHU",
            expected: []string{"charmander", "bulbasaur", "pikachu"},
        },
    }

    for _, c := range cases {
        actual := CleanInput(c.input)
        length := len(actual)
        expectedLength := len(c.expected)
        if length != expectedLength {
            t.Errorf("Slice lengths do not match!")
        }
        for i:= range actual {
            word := actual[i]
            expectedWord := c.expected[i]
            if word != expectedWord {
                t.Errorf("Words do not match!")
            }
        }
    }
}
