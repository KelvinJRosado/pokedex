package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func GetLocationArea(id int) (LocationArea, error) {

	// Build exact path to get location area data
	fullPath := fmt.Sprintf("%v%v%d", POKEAPI_BASE_URL, "location-area/", id)

	// Built GET request
	res, err := http.Get(fullPath)
	if err != nil {
		fmt.Printf("Error calling LocationArea API: %v", err.Error())
		return LocationArea{}, err
	}
	// Read and parse response
	data, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		fmt.Printf("LocationArea API call failed with status code: %d and\ndata: %s\n", res.StatusCode, data)
		return LocationArea{}, errors.New("LocationArea API call failed")
	}
	if err != nil {
		fmt.Printf("Failed to parse response body: %v", err.Error())
		return LocationArea{}, errors.New("LocationArea API response could not be parsed")
	}

	// Convert data to struct
	var location LocationArea
	err = json.Unmarshal(data, &location)
	if err != nil {
		fmt.Printf("Failed to unmarshal data: %v", err.Error())
		return LocationArea{}, errors.New("LocationArea API response could not be parsed")
	}

	return location, nil
}
