package repl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/kelvinjrosado/pokedex/internal/logger"
	"github.com/kelvinjrosado/pokedex/internal/pokeapi"
	"github.com/kelvinjrosado/pokedex/internal/pokecache"
)

type Config struct {
	Client           *pokeapi.Client
	CaughtPokemonMap *pokeapi.CaughtPokemonMap
	MapIndex         int
	Writer           io.Writer
}

// Run reads commands from in and writes prompts and output to out, returning
// when the user runs `exit`, when in reaches EOF, or when the scanner errors.
func Run(in io.Reader, out io.Writer, lgr *logger.CustomLogger) {

	scanner := bufio.NewScanner(in)

	cache := pokecache.NewCache(time.Second * 5)
	defer cache.Stop()

	config := Config{
		Client:           pokeapi.NewClient(cache, lgr),
		CaughtPokemonMap: pokeapi.NewCaughtPokemonMap(),
		Writer:           out,
	}

	lgr.Debug("App started")
	fmt.Fprint(out, "Pokedex > ")
	for scanner.Scan() {
		cleaned := cleanInput(scanner.Text())

		if len(cleaned) > 0 {
			command, ok := getCommand(cleaned[0])
			if !ok {
				lgr.Error("Unknown command", "invalidCommand", cleaned[0])
				fmt.Fprint(out, "Pokedex > ")
				continue
			}

			err := command.callback(&config, cleaned)
			if err != nil {
				if errors.Is(err, ErrCleanExit) {
					return
				}
				lgr.Error("Error executing command", "command", command, "error", err.Error())
			}
		}

		fmt.Fprint(out, "Pokedex > ")
	}

	if err := scanner.Err(); err != nil {
		lgr.Error("Invalid input", "error", err.Error())
	}
}
