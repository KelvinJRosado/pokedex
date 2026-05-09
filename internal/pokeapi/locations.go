package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/kelvinjrosado/pokedex/internal/pokecache"
)

func GetLocationAreaSlice(startId, count int, cache *pokecache.Cache) (LocationAreaList, error) {

	// Build exact path to get location area data
	fullPath := fmt.Sprintf("%vlocation-area/?offset=%d&limit=%d", POKEAPI_BASE_URL, startId, count)

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
			return LocationAreaList{}, err
		}
		defer res.Body.Close()

		// Read and parse response
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return LocationAreaList{}, err
		}

		if res.StatusCode > 299 {
			return LocationAreaList{}, errors.New("LocationArea API call failed")
		}

		// Save response to cache
		cache.Add(fullPath, data)
	}

	// Convert data to struct
	var locationsData LocationAreaList
	err := json.Unmarshal(data, &locationsData)
	if err != nil {
		return LocationAreaList{}, err
	}

	return locationsData, nil
}

func GetLocationAreaDetails(name string, cache *pokecache.Cache) (LocationAreaDetails, error) {

	// Build exact path to get location area detail data
	fullPath := fmt.Sprintf("%vlocation-area/%v", POKEAPI_BASE_URL, name)

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
			return LocationAreaDetails{}, err
		}
		defer res.Body.Close()

		// Read and parse response
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return LocationAreaDetails{}, err
		}

		if res.StatusCode > 299 {
			fmt.Printf("LocationArea API call failed with status code: %d and\ndata: %s\n", res.StatusCode, data)
			return LocationAreaDetails{}, errors.New("LocationArea API call failed")
		}

		// Save response to cache
		cache.Add(fullPath, data)
	}

	// Convert data to struct
	var locationsData LocationAreaDetails
	err := json.Unmarshal(data, &locationsData)
	if err != nil {
		return LocationAreaDetails{}, err
	}

	return locationsData, nil
}
