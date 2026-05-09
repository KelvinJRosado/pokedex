package pokeapi

// pokeapiBaseUrl is a var rather than a const so tests can point it at an
// httptest.Server. Production code never reassigns it.
var pokeapiBaseUrl = "https://pokeapi.co/api/v2/"

const MapIncrement = 20
