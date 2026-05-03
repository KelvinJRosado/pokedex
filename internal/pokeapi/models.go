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

type StatMeta struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type TypeMeta struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokemonType struct {
	Slot int      `json:"slot"`
	Type TypeMeta `json:"type"`
}

type PokemonStat struct {
	BaseStat int      `json:"base_stat"`
	Effort   int      `json:"effort"`
	Stat     StatMeta `json:"stat"`
}

type PokemonDetails struct {
	Id             int           `json:"id"`
	BaseExperience int           `json:"base_experience"`
	Height         int           `json:"height"`
	Weight         int           `json:"weight"`
	Name           string        `json:"name"`
	Stats          []PokemonStat `json:"stats"`
	Types          []PokemonType `json:"types"`
}
