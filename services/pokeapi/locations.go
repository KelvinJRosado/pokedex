package pokeapi

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func GetLocationArea(id int) (LocationArea, error) {

	fullPath := fmt.Sprintf("%v%v%d", POKEAPI_BASE_URL, "location-area/", id)

	res, err := http.Get(fullPath)
	if err != nil {
		fmt.Printf("Error calling LocationArea API: %v", err.Error())
		return LocationArea{}, err
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		fmt.Printf("LocationArea API call failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		return LocationArea{}, errors.New("LocationArea API call failed")
	}
	if err != nil {
		fmt.Printf("Failed to parse response body")
		return LocationArea{}, errors.New("LocationArea API response could not be parsed")
	}
	// TODO: Parse response
	fmt.Printf("%s", body)

	return LocationArea{}, nil
}
