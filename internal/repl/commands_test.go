package repl

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/kelvinjrosado/pokedex/internal/logger"
	"github.com/kelvinjrosado/pokedex/internal/pokeapi"
	"github.com/kelvinjrosado/pokedex/internal/pokecache"
)

// newTestConfig builds a Config wired to an in-memory buffer so tests can
// assert on output without touching os.Stdout.
func newTestConfig(t *testing.T) (*Config, *bytes.Buffer) {
	t.Helper()
	cache := pokecache.NewCache(time.Hour)
	t.Cleanup(cache.Stop)
	lgr := logger.NewWithWriters(io.Discard, io.Discard)
	var buf bytes.Buffer
	return &Config{
		Client:           pokeapi.NewClient(cache, lgr),
		CaughtPokemonMap: pokeapi.NewCaughtPokemonMap(),
		Writer:           &buf,
	}, &buf
}

// withTestAPI points pokeapi at an httptest.Server for the duration of the
// calling test and restores the original base URL via t.Cleanup.
func withTestAPI(t *testing.T, h http.HandlerFunc) {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	t.Cleanup(pokeapi.SetBaseURLForTest(srv.URL + "/"))
}

func TestCommandExit(t *testing.T) {
	cfg, out := newTestConfig(t)
	err := commandExit(cfg, nil)
	if !errors.Is(err, ErrCleanExit) {
		t.Fatalf("expected ErrCleanExit, got %v", err)
	}
	if !strings.Contains(out.String(), "Closing the Pokedex") {
		t.Fatalf("expected farewell message, got: %q", out.String())
	}
}

func TestCommandHelpListsAllCommandsSorted(t *testing.T) {
	cfg, out := newTestConfig(t)
	if err := commandHelp(cfg, nil); err != nil {
		t.Fatalf("commandHelp returned error: %v", err)
	}
	got := out.String()
	for _, cmd := range getAllCommands() {
		if !strings.Contains(got, cmd.name+":") {
			t.Errorf("expected help output to mention command %q; got:\n%s", cmd.name, got)
		}
	}
	idxCatch := strings.Index(got, "catch:")
	idxExit := strings.Index(got, "exit:")
	idxMap := strings.Index(got, "map:")
	if !(idxCatch < idxExit && idxExit < idxMap) {
		t.Errorf("expected commands in alphabetical order, got positions catch=%d exit=%d map=%d", idxCatch, idxExit, idxMap)
	}
}

func TestGetCommandKnown(t *testing.T) {
	for _, name := range []string{"exit", "help", "map", "mapb", "explore", "catch", "inspect", "pokedex"} {
		cmd, ok := getCommand(name)
		if !ok {
			t.Errorf("expected command %q to exist", name)
			continue
		}
		if cmd.name != name {
			t.Errorf("expected name %q, got %q", name, cmd.name)
		}
		if cmd.callback == nil {
			t.Errorf("expected callback for %q", name)
		}
	}
}

func TestGetCommandUnknown(t *testing.T) {
	if _, ok := getCommand("nonsense"); ok {
		t.Fatal("expected unknown command lookup to fail")
	}
}

func TestCommandMapAdvancesIndexAndPrintsResults(t *testing.T) {
	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"results":[{"name":"area-a","url":"x"},{"name":"area-b","url":"y"}]}`))
	})

	cfg, out := newTestConfig(t)
	if err := commandMap(cfg, nil); err != nil {
		t.Fatalf("commandMap returned error: %v", err)
	}
	if !strings.Contains(out.String(), "area-a") || !strings.Contains(out.String(), "area-b") {
		t.Errorf("expected location names in output, got: %q", out.String())
	}
	if cfg.MapIndex != pokeapi.MapIncrement {
		t.Errorf("expected MapIndex to advance to %d, got %d", pokeapi.MapIncrement, cfg.MapIndex)
	}
}

func TestCommandMapPropagatesError(t *testing.T) {
	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	})

	cfg, _ := newTestConfig(t)
	if err := commandMap(cfg, nil); err == nil {
		t.Fatal("expected error from commandMap on 500")
	}
	if cfg.MapIndex != 0 {
		t.Errorf("expected MapIndex to remain 0 on error, got %d", cfg.MapIndex)
	}
}

func TestCommandMapbFirstPageReturnsError(t *testing.T) {
	cfg, _ := newTestConfig(t)
	cfg.MapIndex = pokeapi.MapIncrement // simulate having only fetched the first page
	err := commandMapb(cfg, nil)
	if err == nil || !strings.Contains(err.Error(), "first page") {
		t.Fatalf("expected first-page error, got %v", err)
	}
}

func TestCommandMapbStepsBack(t *testing.T) {
	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"results":[{"name":"prev-area","url":"x"}]}`))
	})

	cfg, out := newTestConfig(t)
	cfg.MapIndex = pokeapi.MapIncrement * 3 // user has paged forward
	if err := commandMapb(cfg, nil); err != nil {
		t.Fatalf("commandMapb returned error: %v", err)
	}
	if !strings.Contains(out.String(), "prev-area") {
		t.Errorf("expected prev-area in output, got: %q", out.String())
	}
	// Net effect of mapb is: index moves back by one page (MapIncrement).
	if want := pokeapi.MapIncrement * 2; cfg.MapIndex != want {
		t.Errorf("expected MapIndex %d after mapb, got %d", want, cfg.MapIndex)
	}
}

func TestCommandMapbRestoresIndexOnError(t *testing.T) {
	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	})

	cfg, _ := newTestConfig(t)
	cfg.MapIndex = pokeapi.MapIncrement * 3
	if err := commandMapb(cfg, nil); err == nil {
		t.Fatal("expected error")
	}
	if want := pokeapi.MapIncrement * 3; cfg.MapIndex != want {
		t.Errorf("expected MapIndex restored to %d, got %d", want, cfg.MapIndex)
	}
}

func TestCommandRequiresArg(t *testing.T) {
	cases := map[string]func(*Config, []string) error{
		"explore": commandExplore,
		"catch":   commandCatch,
		"inspect": commandInspect,
	}
	for name, fn := range cases {
		t.Run(name, func(t *testing.T) {
			cfg, _ := newTestConfig(t)
			if err := fn(cfg, []string{name}); err == nil {
				t.Fatalf("expected error when no argument is provided to %s", name)
			}
		})
	}
}

func TestCommandExploreListsEncounters(t *testing.T) {
	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{
			"id": 1,
			"name": "test-area",
			"pokemon_encounters": [
				{"pokemon": {"name": "rattata", "url": "x"}},
				{"pokemon": {"name": "pidgey", "url": "y"}}
			]
		}`))
	})

	cfg, out := newTestConfig(t)
	if err := commandExplore(cfg, []string{"explore", "test-area"}); err != nil {
		t.Fatalf("commandExplore returned error: %v", err)
	}
	if !strings.Contains(out.String(), "rattata") || !strings.Contains(out.String(), "pidgey") {
		t.Errorf("expected encounters in output, got: %q", out.String())
	}
}

func TestCommandExploreAPIError(t *testing.T) {
	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	})

	cfg, _ := newTestConfig(t)
	if err := commandExplore(cfg, []string{"explore", "bad-area"}); err == nil {
		t.Fatal("expected error from commandExplore when API fails")
	}
}

func TestCommandCatchSuccess(t *testing.T) {
	// BaseExperience 0 ⇒ rand.IntN(maxCatchRate) >= 0 is always true ⇒ always caught.
	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"name":"weedle","base_experience":0}`))
	})

	cfg, out := newTestConfig(t)
	if err := commandCatch(cfg, []string{"catch", "weedle"}); err != nil {
		t.Fatalf("commandCatch returned error: %v", err)
	}
	if !strings.Contains(out.String(), "was caught") {
		t.Errorf("expected catch success message, got: %q", out.String())
	}
	if _, ok := cfg.CaughtPokemonMap.Get("weedle"); !ok {
		t.Error("expected weedle to be added to caught map")
	}
}

func TestCommandCatchEscape(t *testing.T) {
	// BaseExperience == maxCatchRate ⇒ rand.IntN(maxCatchRate) is always < it ⇒ always escapes.
	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"name":"mewtwo","base_experience":700}`))
	})

	cfg, out := newTestConfig(t)
	if err := commandCatch(cfg, []string{"catch", "mewtwo"}); err != nil {
		t.Fatalf("commandCatch returned error: %v", err)
	}
	if !strings.Contains(out.String(), "escaped") {
		t.Errorf("expected escape message, got: %q", out.String())
	}
	if _, ok := cfg.CaughtPokemonMap.Get("mewtwo"); ok {
		t.Error("expected mewtwo to NOT be added to caught map")
	}
}

func TestCommandCatchAPIError(t *testing.T) {
	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})

	cfg, _ := newTestConfig(t)
	if err := commandCatch(cfg, []string{"catch", "missingno"}); err == nil {
		t.Fatal("expected error from commandCatch when API fails")
	}
}

func TestCommandInspectMissNotifiesUser(t *testing.T) {
	cfg, out := newTestConfig(t)
	if err := commandInspect(cfg, []string{"inspect", "missingno"}); err != nil {
		t.Fatalf("commandInspect should not return error on miss, got: %v", err)
	}
	if !strings.Contains(out.String(), "have not caught") {
		t.Errorf("expected miss message, got: %q", out.String())
	}
}

func TestCommandInspectPrintsDetails(t *testing.T) {
	cfg, out := newTestConfig(t)
	cfg.CaughtPokemonMap.Add("eevee", pokeapi.PokemonDetails{
		Name:   "eevee",
		Height: 3,
		Weight: 65,
		Stats: []pokeapi.PokemonStat{
			{BaseStat: 55, Stat: pokeapi.StatMeta{Name: "hp"}},
		},
		Types: []pokeapi.PokemonType{
			{Slot: 1, Type: pokeapi.TypeMeta{Name: "normal"}},
		},
	})

	if err := commandInspect(cfg, []string{"inspect", "eevee"}); err != nil {
		t.Fatalf("commandInspect returned error: %v", err)
	}
	for _, want := range []string{"eevee", "hp", "normal", "55"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("expected output to contain %q, got: %q", want, out.String())
		}
	}
}

func TestCommandPokedexEmpty(t *testing.T) {
	cfg, _ := newTestConfig(t)
	if err := commandPokedex(cfg, nil); err == nil {
		t.Fatal("expected error when pokedex is empty")
	}
}

func TestCommandPokedexListsCaughtSorted(t *testing.T) {
	cfg, out := newTestConfig(t)
	cfg.CaughtPokemonMap.Add("zubat", pokeapi.PokemonDetails{Name: "zubat"})
	cfg.CaughtPokemonMap.Add("abra", pokeapi.PokemonDetails{Name: "abra"})
	cfg.CaughtPokemonMap.Add("magikarp", pokeapi.PokemonDetails{Name: "magikarp"})

	if err := commandPokedex(cfg, nil); err != nil {
		t.Fatalf("commandPokedex returned error: %v", err)
	}

	got := out.String()
	idxAbra := strings.Index(got, "abra")
	idxMagikarp := strings.Index(got, "magikarp")
	idxZubat := strings.Index(got, "zubat")
	if idxAbra < 0 || idxMagikarp < 0 || idxZubat < 0 {
		t.Fatalf("expected all three names in output, got: %q", got)
	}
	if !(idxAbra < idxMagikarp && idxMagikarp < idxZubat) {
		t.Errorf("expected alphabetical order, got positions abra=%d magikarp=%d zubat=%d", idxAbra, idxMagikarp, idxZubat)
	}
}
