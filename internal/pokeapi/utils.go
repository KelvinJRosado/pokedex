package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// fetchWithCache is a free function (not a method) because Go does not allow
// methods to declare type parameters. Callers pass the Client they would have
// been a receiver on.
func fetchWithCache[T any](c *Client, fullPath, label string) (T, error) {
	// Init response
	var response T

	// Store raw data from response
	var data []byte

	// Check cache
	cacheRes, hit := c.Cache.Get(fullPath)
	if hit {
		data = cacheRes
		c.Logger.Debug("Serving response from cache", "path", fullPath)
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
		c.Cache.Add(fullPath, data)
	}

	// Convert data to struct
	err := json.Unmarshal(data, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}
