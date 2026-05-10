package repl

import (
	"errors"
	"io"

	"github.com/kelvinjrosado/pokedex/internal/logger"
	"github.com/kelvinjrosado/pokedex/internal/pokeapi"
	"github.com/kelvinjrosado/pokedex/internal/pokecache"
)

type Config struct {
	Cache            *pokecache.Cache
	CaughtPokemonMap *pokeapi.CaughtPokemonMap
	Logger           *logger.CustomLogger
	MapIndex         int
	Writer           io.Writer
}

// Struct defining the format for a CLI command definition
type cliCommand struct {
	name        string
	description string
	callback    func(*Config, []string) error
}

var ErrCleanExit = errors.New("clean exit")
