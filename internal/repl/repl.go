package repl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/kelvinjrosado/pokedex/internal/logger"
	"github.com/kelvinjrosado/pokedex/internal/pokeapi"
	"github.com/kelvinjrosado/pokedex/internal/pokecache"
)

// Run reads commands from in and writes prompts and output to out, returning
// when the user runs `exit`, when in reaches EOF, or when the scanner errors.
func Run(in io.Reader, out io.Writer) {

	// Init logger
	lgr, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("Unable to start logger: %v", err)
	}
	defer lgr.Close()

	scanner := bufio.NewScanner(in)

	// Init cache
	cache := pokecache.NewCache(time.Second * 5)
	defer cache.Stop()

	// Init pokedex
	cpl := pokeapi.NewCaughtPokemonMap()

	config := Config{
		Cache:            cache,
		CaughtPokemonMap: cpl,
		Writer:           out,
		Logger:           lgr,
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
