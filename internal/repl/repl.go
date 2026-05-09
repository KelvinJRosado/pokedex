package repl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/kelvinjrosado/pokedex/internal/pokeapi"
	"github.com/kelvinjrosado/pokedex/internal/pokecache"
)

// Run reads commands from in and writes prompts and output to out, returning
// when the user runs `exit`, when in reaches EOF, or when the scanner errors.
func Run(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	cache := pokecache.NewCache(time.Second * 5)
	defer cache.Stop()

	cpl := pokeapi.NewCaughtPokemonMap()
	config := Config{
		Cache:            cache,
		CaughtPokemonMap: cpl,
		Writer:           out,
	}

	fmt.Fprint(out, "Pokedex > ")
	for scanner.Scan() {
		cleaned := cleanInput(scanner.Text())

		if len(cleaned) > 0 {
			command, ok := getCommand(cleaned[0])
			if !ok {
				fmt.Fprintln(out, "Unknown command")
				fmt.Fprint(out, "Pokedex > ")
				continue
			}

			err := command.callback(&config, cleaned)
			if err != nil {
				if errors.Is(err, ErrCleanExit) {
					return
				}
				fmt.Fprintf(out, "Error: %v\n", err)
			}
		}

		fmt.Fprint(out, "Pokedex > ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(out, "Invalid input: %s\n", err)
	}
}
