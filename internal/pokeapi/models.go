package pokeapi

type LocationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationAreaList struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type PokemonSummary struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokemonEncounter struct {
	Pokemon PokemonSummary `json:"pokemon"`
}

type LocationAreaDetails struct {
	Id         int                `json:"int"`
	Name       string             `json:"string"`
	Encounters []PokemonEncounter `json:"pokemon_encounters"`
}
