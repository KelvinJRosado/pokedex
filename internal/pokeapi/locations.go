package pokeapi

import "fmt"

func (c *Client) GetLocationAreaSlice(startId, count int) (LocationAreaList, error) {
	fullPath := fmt.Sprintf("%vlocation-area/?offset=%d&limit=%d", pokeapiBaseUrl, startId, count)
	return fetchWithCache[LocationAreaList](c, fullPath, "LocationAreaList")
}

func (c *Client) GetLocationAreaDetails(name string) (LocationAreaDetails, error) {
	fullPath := fmt.Sprintf("%vlocation-area/%v", pokeapiBaseUrl, name)
	return fetchWithCache[LocationAreaDetails](c, fullPath, "LocationAreaDetails")
}
