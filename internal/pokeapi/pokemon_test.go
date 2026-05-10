package pokeapi

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestGetPokemonDetails(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{
			"id": 4,
			"base_experience": 62,
			"height": 6,
			"weight": 85,
			"name": "charmander",
			"stats": [
				{"base_stat": 39, "effort": 0, "stat": {"name": "hp", "url": "x"}}
			],
			"types": [
				{"slot": 1, "type": {"name": "fire", "url": "x"}}
			]
		}`))
	}))
	defer srv.Close()
	withTestServer(t, srv)

	client := newTestClient(t)
	got, err := client.GetPokemonDetails("charmander")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/pokemon/charmander" {
		t.Fatalf("unexpected request path: %q", gotPath)
	}
	if got.Name != "charmander" || got.BaseExperience != 62 || got.Height != 6 {
		t.Fatalf("unexpected pokemon: %+v", got)
	}
	if len(got.Stats) != 1 || got.Stats[0].Stat.Name != "hp" {
		t.Fatalf("unexpected stats: %+v", got.Stats)
	}
	if len(got.Types) != 1 || got.Types[0].Type.Name != "fire" {
		t.Fatalf("unexpected types: %+v", got.Types)
	}
}

func TestCaughtPokemonMapAddGet(t *testing.T) {
	m := NewCaughtPokemonMap()

	if _, ok := m.Get("absent"); ok {
		t.Fatal("expected miss on empty map")
	}

	want := PokemonDetails{Name: "pikachu", BaseExperience: 112}
	m.Add("pikachu", want)

	got, ok := m.Get("pikachu")
	if !ok {
		t.Fatal("expected hit after Add")
	}
	if got.Name != want.Name || got.BaseExperience != want.BaseExperience {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestCaughtPokemonMapOverwrites(t *testing.T) {
	m := NewCaughtPokemonMap()
	m.Add("eevee", PokemonDetails{Name: "eevee", BaseExperience: 1})
	m.Add("eevee", PokemonDetails{Name: "eevee", BaseExperience: 65})

	got, ok := m.Get("eevee")
	if !ok || got.BaseExperience != 65 {
		t.Fatalf("expected overwritten value, got %+v ok=%v", got, ok)
	}
}

// Sanity check that concurrent Add/Get is race-free under -race.
func TestCaughtPokemonMapConcurrent(t *testing.T) {
	m := NewCaughtPokemonMap()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			m.Add("x", PokemonDetails{Name: "x"})
		}()
		go func() {
			defer wg.Done()
			_, _ = m.Get("x")
		}()
	}
	wg.Wait()
}

func TestCaughtPokemonMapGetAll(t *testing.T) {
	m := NewCaughtPokemonMap()

	got := m.GetAll()
	if len(got) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(got))
	}

	m.Add("pikachu", PokemonDetails{Name: "pikachu", BaseExperience: 112})
	m.Add("bulbasaur", PokemonDetails{Name: "bulbasaur", BaseExperience: 64})

	got = m.GetAll()
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got["pikachu"].Name != "pikachu" {
		t.Errorf("expected pikachu, got %s", got["pikachu"].Name)
	}
	if got["bulbasaur"].Name != "bulbasaur" {
		t.Errorf("expected bulbasaur, got %s", got["bulbasaur"].Name)
	}

	// Verify GetAll returns a copy (modifications don't affect the original).
	got["charmander"] = PokemonDetails{Name: "charmander"}
	if _, ok := m.Get("charmander"); ok {
		t.Error("modification to GetAll result should not affect original map")
	}
}
