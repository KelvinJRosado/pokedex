package pokeapi

import (
	"encoding/json"
	"errors"
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
			fmt.Printf("Error calling Pokemon details API: %v\n", err.Error())
			return PokemonDetails{}, err
		}
		// Read and parse response
		data, err = io.ReadAll(res.Body)
		defer res.Body.Close()
		if res.StatusCode > 299 {
			fmt.Printf("Pokemon details API call failed with status code: %d and\ndata: %s\n", res.StatusCode, data)
			return PokemonDetails{}, errors.New("Pokemon details API call failed")
		}
		if err != nil {
			fmt.Printf("Failed to parse response body: %v\n", err.Error())
			return PokemonDetails{}, errors.New("Pokemon details API response could not be parsed")
		}

		// Save response to cache
		cache.Add(fullPath, data)
	}

	// Convert data to struct
	var pokemonData PokemonDetails
	err := json.Unmarshal(data, &pokemonData)
	if err != nil {
		fmt.Printf("Failed to unmarshal data: %v\n", err.Error())
		return PokemonDetails{}, errors.New("Pokemon details API response could not be parsed")
	}

	return pokemonData, nil
}

type CaughtPokemonMap struct {
	entries map[string]PokemonDetails
	mu      sync.RWMutex
}

func NewCaughtPokemonMap() *CaughtPokemonMap {
	// Initialize map
	initEntries := make(map[string]PokemonDetails)

	res := CaughtPokemonMap{
		entries: initEntries,
	}

	return &res
}

// Create or update a cache entry
func (c *CaughtPokemonMap) Add(key string, myVal PokemonDetails) {

	// take lock
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = myVal
}

// Retrieve a cache entry
func (c *CaughtPokemonMap) Get(key string) (PokemonDetails, bool) {

	// Take read lock
	c.mu.RLock()
	defer c.mu.RUnlock()

	res, ok := c.entries[key]
	if !ok {
		return PokemonDetails{}, false
	}

	return res, true

}
