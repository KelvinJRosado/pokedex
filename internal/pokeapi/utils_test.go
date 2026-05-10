package pokeapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/kelvinjrosado/pokedex/internal/logger"
	"github.com/kelvinjrosado/pokedex/internal/pokecache"
)

// withTestServer points pokeapiBaseUrl at the given test server for the
// duration of the calling test and restores the original value afterwards.
func withTestServer(t *testing.T, srv *httptest.Server) {
	t.Helper()
	orig := pokeapiBaseUrl
	pokeapiBaseUrl = srv.URL + "/"
	t.Cleanup(func() {
		pokeapiBaseUrl = orig
	})
}

func newCache(t *testing.T) *pokecache.Cache {
	t.Helper()
	// Long interval so reapLoop never runs during the test.
	c := pokecache.NewCache(time.Hour)
	t.Cleanup(c.Stop)
	return c
}

func newTestClient(t *testing.T) *Client {
	t.Helper()
	return NewClient(newCache(t), logger.NewWithWriters(io.Discard, io.Discard))
}

func TestFetchWithCacheNetworkSuccess(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"charmander","base_experience":62}`))
	}))
	defer srv.Close()
	withTestServer(t, srv)

	client := newTestClient(t)
	got, err := fetchWithCache[PokemonDetails](client, pokeapiBaseUrl+"pokemon/charmander", "PokemonDetails")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "charmander" || got.BaseExperience != 62 {
		t.Fatalf("unexpected response: %+v", got)
	}
	if hits != 1 {
		t.Fatalf("expected 1 hit, got %d", hits)
	}
}

func TestFetchWithCacheServesFromCache(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		_, _ = w.Write([]byte(`{"name":"pikachu"}`))
	}))
	defer srv.Close()
	withTestServer(t, srv)

	client := newTestClient(t)
	url := pokeapiBaseUrl + "pokemon/pikachu"

	if _, err := fetchWithCache[PokemonDetails](client, url, "PokemonDetails"); err != nil {
		t.Fatalf("first call failed: %v", err)
	}
	if _, err := fetchWithCache[PokemonDetails](client, url, "PokemonDetails"); err != nil {
		t.Fatalf("second call failed: %v", err)
	}
	if hits != 1 {
		t.Fatalf("expected exactly 1 network hit, got %d", hits)
	}
}

func TestFetchWithCacheNon2xxReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer srv.Close()
	withTestServer(t, srv)

	client := newTestClient(t)
	_, err := fetchWithCache[PokemonDetails](client, pokeapiBaseUrl+"pokemon/missingno", "PokemonDetails")
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Fatalf("expected error to mention 404, got: %v", err)
	}
}

func TestFetchWithCacheMalformedJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`not json`))
	}))
	defer srv.Close()
	withTestServer(t, srv)

	client := newTestClient(t)
	_, err := fetchWithCache[PokemonDetails](client, pokeapiBaseUrl+"pokemon/x", "PokemonDetails")
	if err == nil {
		t.Fatal("expected unmarshal error")
	}
}

func TestFetchWithCacheNetworkError(t *testing.T) {
	// An unroutable URL forces http.Get to fail.
	client := newTestClient(t)
	_, err := fetchWithCache[PokemonDetails](client, "http://127.0.0.1:1/never", "PokemonDetails")
	if err == nil {
		t.Fatal("expected network error")
	}
}

func TestFetchWithCacheReadBodyError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			t.Fatal("server does not support hijacking")
		}
		conn, buf, err := hj.Hijack()
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		// Write a partial response then close, causing io.ReadAll to fail.
		_, _ = buf.WriteString("HTTP/1.1 200 OK\r\nContent-Length: 100\r\n\r\npartial")
		buf.Flush()
		conn.Close()
	}))
	defer srv.Close()
	withTestServer(t, srv)

	client := newTestClient(t)
	_, err := fetchWithCache[PokemonDetails](client, pokeapiBaseUrl+"pokemon/glitcho", "PokemonDetails")
	if err == nil {
		t.Fatal("expected error reading response body")
	}
}
