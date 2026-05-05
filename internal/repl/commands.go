package repl

import (
	"fmt"
	"math/rand"
	"slices"

	"github.com/kelvinjrosado/pokedex/internal/pokeapi"
)

func commandExit(config *Config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	return CleanExit
}

func commandHelp(config *Config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")

	// Sort map
	keys := make([]string, 0, len(commandRegistry))
	for k := range commandRegistry {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	// Print commands in alphabetical order
	for _, k := range keys {
		fmt.Printf("%v: %v\n", k, commandRegistry[k].description)
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

func commandInspect(config *Config, args []string) error {

	// Grab from caught list
	name := args[1]
	val, ok := config.CaughtPokemonMap.Get(name)

	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	// If caught, print info
	fmt.Printf("Name: %v\n", val.Name)
	fmt.Printf("Height: %d\n", val.Height)
	fmt.Printf("Weight: %d\n", val.Weight)
	fmt.Println("Stats:")
	for _, s := range val.Stats {
		fmt.Printf("  -%v: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range val.Types {
		fmt.Printf("  - %v\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(config *Config, args []string) error {

	caught := config.CaughtPokemonMap.Entries

	// Base case: Pokedex is empty
	if len(caught) == 0 {
		fmt.Println("You have not caught any Pokemon yet")
	}

	// Sort map
	keys := make([]string, 0, len(caught))
	for k := range caught {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for k, _ := range caught {
		fmt.Printf(" - %v\n", k)
	}

	return nil
}
