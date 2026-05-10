package repl

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"

	"github.com/kelvinjrosado/pokedex/internal/pokeapi"
)

var ErrCleanExit = errors.New("clean exit")

func commandExit(config *Config, args []string) error {
	fmt.Fprintln(config.Writer, "Closing the Pokedex... Goodbye!")
	return ErrCleanExit
}

func commandHelp(config *Config, args []string) error {
	fmt.Fprintln(config.Writer, "Welcome to the Pokedex!")
	fmt.Fprint(config.Writer, "Usage:\n\n")

	cmds := getAllCommands()
	slices.SortFunc(cmds, func(a, b cliCommand) int {
		return strings.Compare(a.name, b.name)
	})

	for _, cmd := range cmds {
		fmt.Fprintf(config.Writer, "%v: %v\n", cmd.name, cmd.description)
	}

	return nil
}

func commandMap(config *Config, args []string) error {

	las, err := config.Client.GetLocationAreaSlice(config.MapIndex, pokeapi.MapIncrement)
	if err != nil {
		return err
	}

	for _, la := range las.Results {
		fmt.Fprintln(config.Writer, la.Name)
	}

	// Increment map pointer
	config.MapIndex += pokeapi.MapIncrement

	return nil
}

func commandMapb(config *Config, args []string) error {
	// Check base case
	if config.MapIndex <= pokeapi.MapIncrement {
		return errors.New("you're on the first page")
	}

	// Decrease map pointer
	config.MapIndex -= (pokeapi.MapIncrement * 2)

	las, err := config.Client.GetLocationAreaSlice(config.MapIndex, pokeapi.MapIncrement)
	if err != nil {
		config.MapIndex += (pokeapi.MapIncrement * 2) // restore map index
		return err
	}

	for _, la := range las.Results {
		fmt.Fprintln(config.Writer, la.Name)
	}

	// Increase map pointer again as we travelled
	config.MapIndex += pokeapi.MapIncrement

	return nil
}

func commandExplore(config *Config, args []string) error {

	// Check for args being present
	if len(args) < 2 {
		return errors.New("insufficient args provided for \"explore\"")
	}

	locationName := args[1]

	details, err := config.Client.GetLocationAreaDetails(locationName)
	if err != nil {
		return err
	}

	fmt.Fprintf(config.Writer, "Exploring %v...\n", locationName)
	fmt.Fprintln(config.Writer, "Found Pokemon:")

	for _, v := range details.Encounters {
		fmt.Fprintf(config.Writer, " - %v\n", v.Pokemon.Name)
	}

	return nil
}

func commandCatch(config *Config, args []string) error {

	// Check for args being present
	if len(args) < 2 {
		return errors.New("insufficient args provided for \"catch\"")
	}

	pokemonName := args[1]

	// Get details
	pokemonDetails, err := config.Client.GetPokemonDetails(pokemonName)
	if err != nil {
		return err
	}

	fmt.Fprintf(config.Writer, "Throwing a Pokeball at %v...\n", pokemonDetails.Name)

	// Check if caught
	roll := rand.IntN(maxCatchRate)     // Get random number from 0 to max
	be := pokemonDetails.BaseExperience // Get Pokemon base experience for roll

	if roll >= be {
		// Caught
		fmt.Fprintf(config.Writer, "%v was caught!\n", pokemonDetails.Name)

		// Save to caught list
		config.CaughtPokemonMap.Add(pokemonDetails.Name, pokemonDetails)
	} else {
		fmt.Fprintf(config.Writer, "%v escaped!\n", pokemonDetails.Name)
	}

	return nil
}

func commandInspect(config *Config, args []string) error {

	// Check for args being present
	if len(args) < 2 {
		return errors.New("insufficient args provided for \"inspect\"")
	}

	// Grab from caught list
	name := args[1]
	val, ok := config.CaughtPokemonMap.Get(name)

	if !ok {
		fmt.Fprintln(config.Writer, "you have not caught that pokemon")
		return nil
	}

	// If caught, print info
	fmt.Fprintf(config.Writer, "Name: %v\n", val.Name)
	fmt.Fprintf(config.Writer, "Height: %d\n", val.Height)
	fmt.Fprintf(config.Writer, "Weight: %d\n", val.Weight)
	fmt.Fprintln(config.Writer, "Stats:")
	for _, s := range val.Stats {
		fmt.Fprintf(config.Writer, "  -%v: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Fprintln(config.Writer, "Types:")
	for _, t := range val.Types {
		fmt.Fprintf(config.Writer, "  - %v\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(config *Config, args []string) error {

	caught := config.CaughtPokemonMap.GetAll()

	// Base case: Pokedex is empty
	if len(caught) == 0 {
		return errors.New("you have not caught any Pokemon yet")
	}

	// Sort map
	keys := make([]string, 0, len(caught))
	for k := range caught {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for _, k := range keys {
		fmt.Fprintf(config.Writer, " - %v\n", k)
	}

	return nil
}
