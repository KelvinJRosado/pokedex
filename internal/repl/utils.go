package repl

import (
	"errors"
	"fmt"
	"math/rand"
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
	callback    func(*Config, []string) error
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
		"explore": {
			name:        "explore",
			description: "Displays the name of Pokemon that can be encountered in the specified area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch the specified Pokemon",
			callback:    commandCatch,
		},
	}
}

var CleanExit = errors.New("Clean exit")

func commandExit(config *Config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	return CleanExit
}

func commandHelp(config *Config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")

	for k, v := range commandRegistry {
		fmt.Printf("%v: %v\n", k, v.description)
	}

	return nil
}

// Keep track of current map pointer
var mapIndex = 0

func commandMap(config *Config, args []string) error {

	las, err := pokeapi.GetLocationAreaSlice(mapIndex, pokeapi.MAP_INCREMENT, config.Cache)
	if err != nil {
		fmt.Printf("Failed to get location area info: %v\n", err.Error())
		return err
	}

	for _, la := range las.Results {
		fmt.Println(la.Name)
	}

	// Increment map pointer
	mapIndex += pokeapi.MAP_INCREMENT

	return nil
}

func commandMapb(config *Config, args []string) error {
	// Check base case
	if mapIndex <= pokeapi.MAP_INCREMENT {
		fmt.Println("you're on the first page")
		return nil
	}

	// Decrease map pointer
	mapIndex -= (pokeapi.MAP_INCREMENT * 2)

	las, err := pokeapi.GetLocationAreaSlice(mapIndex, pokeapi.MAP_INCREMENT, config.Cache)
	if err != nil {
		fmt.Printf("Failed to get location area info: %v\n", err.Error())
		mapIndex += (pokeapi.MAP_INCREMENT * 2) // restore map index
		return err
	}

	for _, la := range las.Results {
		fmt.Println(la.Name)
	}

	// Increase map pointer again as we travelled
	// 	// Decrease map pointer
	mapIndex += pokeapi.MAP_INCREMENT

	return nil
}

func commandExplore(config *Config, args []string) error {
	locationName := args[1]

	details, err := pokeapi.GetLocationAreaDetails(locationName, config.Cache)
	if err != nil {
		fmt.Printf("Failed to get location area details: %v\n", err.Error())
		return err
	}

	fmt.Printf("Exploring %v...\n", locationName)
	fmt.Println("Found Pokemon:")

	for _, v := range details.Encounters {
		fmt.Printf(" - %v\n", v.Pokemon.Name)
	}

	return nil
}

func commandCatch(config *Config, args []string) error {
	pokemonName := args[1]

	// Get details
	pokemonDetails, err := pokeapi.GetPokemonDetails(pokemonName, config.Cache)
	if err != nil {
		fmt.Printf("Failed to get pokemon details: %v\n", err.Error())
		return err
	}

	fmt.Printf("Throwing a Pokeball at %v...\n", pokemonDetails.Name)

	// Check if caught
	roll := rand.Intn(700)              // Get random number from 0 to max
	be := pokemonDetails.BaseExperience // Get Pokemon base experience for roll

	if roll >= be {
		// Caught
		fmt.Printf("%v was caught!\n", pokemonDetails.Name)

		// Save to caught list
		config.CaughtPokemonMap.Add(pokemonDetails.Name, pokemonDetails)
	} else {
		fmt.Printf("%v escaped!\n", pokemonDetails.Name)
	}

	return nil
}
