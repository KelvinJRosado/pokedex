package repl

import (
	"errors"

	"github.com/kelvinjrosado/pokedex/internal/pokeapi"
	"github.com/kelvinjrosado/pokedex/internal/pokecache"
)

type Config struct {
	Cache            *pokecache.Cache
	CaughtPokemonMap *pokeapi.CaughtPokemonMap
	MapIndex         int
}

// Struct defining the format for a CLI command definition
type cliCommand struct {
	name        string
	description string
	callback    func(*Config, []string) error
}

var CleanExit = errors.New("Clean exit")
