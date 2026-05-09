package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/kelvinjrosado/pokedex/internal/pokecache"
)

func GetPokemonDetails(name string, cache *pokecache.Cache) (PokemonDetails, error) {
	// Build exact path to get location area detail data
	fullPath := fmt.Sprintf("%vpokemon/%v", POKEAPI_BASE_URL, name)

	// Store raw data from response
	var data []byte

	// Check cache
	cacheRes, hit := cache.Get(fullPath)
	if hit {
		data = cacheRes
		fmt.Println("Serving response from cache")
	} else {
		// Built GET request
		res, err := http.Get(fullPath)
		if err != nil {
			return PokemonDetails{}, err
		}
		defer res.Body.Close()

		// Read and parse response
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return PokemonDetails{}, err
		}

		if res.StatusCode > 299 {
			return PokemonDetails{}, fmt.Errorf("Pokemon details API call failed with status %d: %s", res.StatusCode, data)
		}

		// Save response to cache
		cache.Add(fullPath, data)
	}

	// Convert data to struct
	var pokemonData PokemonDetails
	err := json.Unmarshal(data, &pokemonData)
	if err != nil {
		return PokemonDetails{}, err
	}

	return pokemonData, nil
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
