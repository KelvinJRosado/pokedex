package repl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kelvinjrosado/pokedex/internal/pokeapi"
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
	callback    func(*Config) error
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
			description: "Displays the name of the 20 next locations in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the name of the 20 previous locations in the Pokemon world",
			callback:    commandMapb,
		},
	}
}

var CleanExit = errors.New("Clean exit")

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	return CleanExit
}

func commandHelp(config *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")

	for k, v := range commandRegistry {
		fmt.Printf("%v: %v\n", k, v.description)
	}

	return nil
}

// Keep track of current map pointer
var mapIndex = 0

func commandMap(config *Config) error {

	las, err := pokeapi.GetLocationAreaSlice(mapIndex, pokeapi.MAP_INCREMENT, config.Cache)
	if err != nil {
		fmt.Printf("Failed to get location area info: %v", err.Error())
		return err
	}

	for _, la := range las.Results {
		fmt.Println(la.Name)
	}

	// Increment map pointer
	mapIndex += pokeapi.MAP_INCREMENT

	return nil
}

func commandMapb(config *Config) error {
	// Check base case
	if mapIndex <= 20 {
		fmt.Println("you're on the first page")
		return nil
	}

	// Decrease map pointer
	mapIndex -= (pokeapi.MAP_INCREMENT * 2)

	las, err := pokeapi.GetLocationAreaSlice(mapIndex, pokeapi.MAP_INCREMENT, config.Cache)
	if err != nil {
		fmt.Printf("Failed to get location area info: %v", err.Error())
		return err
	}

	for _, la := range las.Results {
		fmt.Println(la.Name)
	}

	return nil
}
