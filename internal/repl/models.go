package repl

import (
	"github.com/kelvinjrosado/pokedex/internal/pokeapi"
	"github.com/kelvinjrosado/pokedex/internal/pokecache"
)

type Config struct {
	Cache            *pokecache.Cache
	CaughtPokemonMap *pokeapi.CaughtPokemonMap
}
