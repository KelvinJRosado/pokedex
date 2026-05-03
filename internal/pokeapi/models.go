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
	Id         int                `json:"id"`
	Name       string             `json:"name"`
	Encounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonDetails struct {
	Id             int    `json:"id"`
	BaseExperience int    `json:"base_experience"`
	Name           string `json:"name"`
}
