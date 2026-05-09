package pokeapi

import (
	"fmt"
	"sync"

	"github.com/kelvinjrosado/pokedex/internal/pokecache"
)

func GetPokemonDetails(name string, cache *pokecache.Cache) (PokemonDetails, error) {
	// Build exact path to get pokemon detail data
	fullPath := fmt.Sprintf("%vpokemon/%v", POKEAPI_BASE_URL, name)

	return fetchWithCache[PokemonDetails](cache, fullPath, "PokemonDetails")
}

type CaughtPokemonMap struct {
	Entries map[string]PokemonDetails
	mu      sync.RWMutex
}

func NewCaughtPokemonMap() *CaughtPokemonMap {
	// Initialize map
	initEntries := make(map[string]PokemonDetails)

	res := CaughtPokemonMap{
		Entries: initEntries,
	}

	return &res
}

// Create or update a cache entry
func (c *CaughtPokemonMap) Add(key string, myVal PokemonDetails) {

	// take lock
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Entries[key] = myVal
}

// Retrieve a cache entry
func (c *CaughtPokemonMap) Get(key string) (PokemonDetails, bool) {

	// Take read lock
	c.mu.RLock()
	defer c.mu.RUnlock()

	res, ok := c.Entries[key]
	if !ok {
		return PokemonDetails{}, false
	}

	return res, true

}
