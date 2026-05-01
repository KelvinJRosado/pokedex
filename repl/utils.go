package repl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kelvinjrosado/pokedex/services/pokeapi"
)

// Split the user's input into words based on whitespace.
// It should also lowercase the input and trim any leading or trailing whitespace.
func cleanInput(text string) []string {
	var out []string

	// First we make everything lowercase
	lower := strings.ToLower(text)

	// Lastly we split on whitespace, removing leading/trailing whitespace
	out = strings.Fields(lower)

	return out
}

// Struct defining the format for a CLI command definition
type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commandRegistry = map[string]cliCommand{}

func initRegistry() {
	commandRegistry = map[string]cliCommand{
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
			description: "Displays the name of locations in the Pokemon world",
			callback:    commandMap,
		},
	}
}

var CleanExit = errors.New("Clean exit")

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	return CleanExit
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")

	for k, v := range commandRegistry {
		fmt.Printf("%v: %v\n", k, v.description)
	}

	return nil
}

func commandMap() error {

	la, err := pokeapi.GetLocationArea(1)
	if err != nil {
		fmt.Printf("Failed to get location area info: %v", err.Error())
		return err
	}

	fmt.Println(la.Name)

	return nil
}
