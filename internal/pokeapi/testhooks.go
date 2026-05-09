package pokeapi

// SetBaseURLForTest swaps the package-level base URL and returns a function
// that restores the previous value. Intended for use from external test
// packages (e.g., the repl package's tests) that need to redirect API calls
// at an httptest.Server. Production code must not call this.
func SetBaseURLForTest(url string) func() {
	prev := pokeapiBaseUrl
	pokeapiBaseUrl = url
	return func() { pokeapiBaseUrl = prev }
}
