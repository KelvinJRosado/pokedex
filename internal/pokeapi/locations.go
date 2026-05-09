package pokeapi

import (
	"fmt"

	"github.com/kelvinjrosado/pokedex/internal/pokecache"
)

func GetLocationAreaSlice(startId, count int, cache *pokecache.Cache) (LocationAreaList, error) {

	// Build exact path to get location area data
	fullPath := fmt.Sprintf("%vlocation-area/?offset=%d&limit=%d", pokeapiBaseUrl, startId, count)

	return fetchWithCache[LocationAreaList](cache, fullPath, "LocationAreaList")
}

func GetLocationAreaDetails(name string, cache *pokecache.Cache) (LocationAreaDetails, error) {

	// Build exact path to get location area detail data
	fullPath := fmt.Sprintf("%vlocation-area/%v", pokeapiBaseUrl, name)

	return fetchWithCache[LocationAreaDetails](cache, fullPath, "LocationAreaDetails")
}
