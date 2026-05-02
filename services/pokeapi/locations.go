package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func GetLocationAreaSlice(startId, count int) (LocationAreaList, error) {

	// Build exact path to get location area data
	fullPath := fmt.Sprintf("%vlocation-area/?offset=%d&limit=%d", POKEAPI_BASE_URL, startId, count)

	// Built GET request
	res, err := http.Get(fullPath)
	if err != nil {
		fmt.Printf("Error calling LocationArea API: %v\n", err.Error())
		return LocationAreaList{}, err
	}
	// Read and parse response
	data, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		fmt.Printf("LocationArea API call failed with status code: %d and\ndata: %s\n", res.StatusCode, data)
		return LocationAreaList{}, errors.New("LocationArea API call failed")
	}
	if err != nil {
		fmt.Printf("Failed to parse response body: %v\n", err.Error())
		return LocationAreaList{}, errors.New("LocationArea API response could not be parsed")
	}

	// Convert data to struct
	var locationsData LocationAreaList
	err = json.Unmarshal(data, &locationsData)
	if err != nil {
		fmt.Printf("Failed to unmarshal data: %v\n", err.Error())
		return LocationAreaList{}, errors.New("LocationArea API response could not be parsed")
	}

	return locationsData, nil
}
