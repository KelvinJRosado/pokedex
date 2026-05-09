package pokeapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLocationAreaSlice(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path + "?" + r.URL.RawQuery
		_, _ = w.Write([]byte(`{
			"count": 2,
			"next": "",
			"previous": "",
			"results": [
				{"name": "canalave-city-area", "url": "x"},
				{"name": "eterna-city-area", "url": "y"}
			]
		}`))
	}))
	defer srv.Close()
	withTestServer(t, srv)

	cache := newCache(t)
	got, err := GetLocationAreaSlice(40, MapIncrement, cache)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "/location-area/?offset=40&limit=20"; gotPath != want {
		t.Fatalf("expected request path %q, got %q", want, gotPath)
	}
	if len(got.Results) != 2 || got.Results[0].Name != "canalave-city-area" {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestGetLocationAreaSlicePropagatesError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()
	withTestServer(t, srv)

	cache := newCache(t)
	if _, err := GetLocationAreaSlice(0, MapIncrement, cache); err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestGetLocationAreaDetails(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{
			"id": 1,
			"name": "canalave-city-area",
			"pokemon_encounters": [
				{"pokemon": {"name": "tentacool", "url": "x"}},
				{"pokemon": {"name": "tentacruel", "url": "y"}}
			]
		}`))
	}))
	defer srv.Close()
	withTestServer(t, srv)

	cache := newCache(t)
	got, err := GetLocationAreaDetails("canalave-city-area", cache)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/location-area/canalave-city-area" {
		t.Fatalf("unexpected request path: %q", gotPath)
	}
	if len(got.Encounters) != 2 || got.Encounters[1].Pokemon.Name != "tentacruel" {
		t.Fatalf("unexpected response: %+v", got)
	}
}
