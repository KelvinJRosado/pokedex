package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kelvinjrosado/pokedex/internal/pokecache"
)

func fetchWithCache[T any](cache *pokecache.Cache, fullPath, label string) (T, error) {
	// Init response
	var response T

	// Store raw data from response
	var data []byte

	// Check cache
	cacheRes, hit := cache.Get(fullPath)
	if hit {
		data = cacheRes
		fmt.Println("Serving response from cache") // TODO: Convert to debug log
	} else {
		// Built GET request
		res, err := http.Get(fullPath)
		if err != nil {
			return response, err
		}
		defer res.Body.Close()

		// Read and parse response. Note: assignment with `=` (not `:=`) so we
		// write to the outer `data`; otherwise it stays nil and Unmarshal fails.
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return response, err
		}

		if res.StatusCode > 299 {
			return response, fmt.Errorf("%v API call failed with status %d: %s", label, res.StatusCode, data)
		}

		// Save response to cache
		cache.Add(fullPath, data)
	}

	// Convert data to struct
	err := json.Unmarshal(data, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}
