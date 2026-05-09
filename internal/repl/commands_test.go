package repl

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kelvinjrosado/pokedex/internal/pokeapi"
	"github.com/kelvinjrosado/pokedex/internal/pokecache"
)

// captureStdout runs fn with os.Stdout redirected to a pipe and returns the
// captured output. Tests use it to assert on user-facing prints without
// coupling to internal formatting helpers.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	done := make(chan string, 1)
	go func() {
		buf, _ := io.ReadAll(r)
		done <- string(buf)
	}()

	fn()
	_ = w.Close()
	return <-done
}

func newTestConfig(t *testing.T) *Config {
	t.Helper()
	cache := pokecache.NewCache(time.Hour)
	t.Cleanup(cache.Stop)
	return &Config{
		Cache:            cache,
		CaughtPokemonMap: pokeapi.NewCaughtPokemonMap(),
	}
}

func TestCommandExit(t *testing.T) {
	cfg := newTestConfig(t)
	out := captureStdout(t, func() {
		err := commandExit(cfg, nil)
		if !errors.Is(err, ErrCleanExit) {
			t.Fatalf("expected ErrCleanExit, got %v", err)
		}
	})
	if !strings.Contains(out, "Closing the Pokedex") {
		t.Fatalf("expected farewell message, got: %q", out)
	}
}

func TestCommandHelpListsAllCommandsSorted(t *testing.T) {
	cfg := newTestConfig(t)
	out := captureStdout(t, func() {
		if err := commandHelp(cfg, nil); err != nil {
			t.Fatalf("commandHelp returned error: %v", err)
		}
	})
	// Every registered command name must appear.
	for _, cmd := range getAllCommands() {
		if !strings.Contains(out, cmd.name+":") {
			t.Errorf("expected help output to mention command %q; got:\n%s", cmd.name, out)
		}
	}
	// Output should be alphabetically sorted by command name.
	idxCatch := strings.Index(out, "catch:")
	idxExit := strings.Index(out, "exit:")
	idxMap := strings.Index(out, "map:")
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
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"results":[{"name":"area-a","url":"x"},{"name":"area-b","url":"y"}]}`))
	}))
	defer srv.Close()
	defer pokeapi.SetBaseURLForTest(srv.URL + "/")()

	cfg := newTestConfig(t)
	out := captureStdout(t, func() {
		if err := commandMap(cfg, nil); err != nil {
			t.Fatalf("commandMap returned error: %v", err)
		}
	})
	if !strings.Contains(out, "area-a") || !strings.Contains(out, "area-b") {
		t.Errorf("expected location names in output, got: %q", out)
	}
	if cfg.MapIndex != pokeapi.MapIncrement {
		t.Errorf("expected MapIndex to advance to %d, got %d", pokeapi.MapIncrement, cfg.MapIndex)
	}
}

func TestCommandMapPropagatesError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()
	defer pokeapi.SetBaseURLForTest(srv.URL + "/")()

	cfg := newTestConfig(t)
	if err := commandMap(cfg, nil); err == nil {
		t.Fatal("expected error from commandMap on 500")
	}
	if cfg.MapIndex != 0 {
		t.Errorf("expected MapIndex to remain 0 on error, got %d", cfg.MapIndex)
	}
}

func TestCommandMapbFirstPageReturnsError(t *testing.T) {
	cfg := newTestConfig(t)
	cfg.MapIndex = pokeapi.MapIncrement // simulate having only fetched the first page
	err := commandMapb(cfg, nil)
	if err == nil || !strings.Contains(err.Error(), "first page") {
		t.Fatalf("expected first-page error, got %v", err)
	}
}

func TestCommandMapbStepsBack(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"results":[{"name":"prev-area","url":"x"}]}`))
	}))
	defer srv.Close()
	defer pokeapi.SetBaseURLForTest(srv.URL + "/")()

	cfg := newTestConfig(t)
	cfg.MapIndex = pokeapi.MapIncrement * 3 // user has paged forward
	out := captureStdout(t, func() {
		if err := commandMapb(cfg, nil); err != nil {
			t.Fatalf("commandMapb returned error: %v", err)
		}
	})
	if !strings.Contains(out, "prev-area") {
		t.Errorf("expected prev-area in output, got: %q", out)
	}
	// Net effect of mapb is: index moves back by one page (MapIncrement).
	if want := pokeapi.MapIncrement * 2; cfg.MapIndex != want {
		t.Errorf("expected MapIndex %d after mapb, got %d", want, cfg.MapIndex)
	}
}

func TestCommandMapbRestoresIndexOnError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()
	defer pokeapi.SetBaseURLForTest(srv.URL + "/")()

	cfg := newTestConfig(t)
	cfg.MapIndex = pokeapi.MapIncrement * 3
	if err := commandMapb(cfg, nil); err == nil {
		t.Fatal("expected error")
	}
	if want := pokeapi.MapIncrement * 3; cfg.MapIndex != want {
		t.Errorf("expected MapIndex restored to %d, got %d", want, cfg.MapIndex)
	}
}

func TestCommandExploreRequiresArg(t *testing.T) {
	cfg := newTestConfig(t)
	if err := commandExplore(cfg, []string{"explore"}); err == nil {
		t.Fatal("expected error when no location name is provided")
	}
}

func TestCommandExploreListsEncounters(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{
			"id": 1,
			"name": "test-area",
			"pokemon_encounters": [
				{"pokemon": {"name": "rattata", "url": "x"}},
				{"pokemon": {"name": "pidgey", "url": "y"}}
			]
		}`))
	}))
	defer srv.Close()
	defer pokeapi.SetBaseURLForTest(srv.URL + "/")()

	cfg := newTestConfig(t)
	out := captureStdout(t, func() {
		if err := commandExplore(cfg, []string{"explore", "test-area"}); err != nil {
			t.Fatalf("commandExplore returned error: %v", err)
		}
	})
	if !strings.Contains(out, "rattata") || !strings.Contains(out, "pidgey") {
		t.Errorf("expected encounters in output, got: %q", out)
	}
}

func TestCommandExploreAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer srv.Close()
	defer pokeapi.SetBaseURLForTest(srv.URL + "/")()

	cfg := newTestConfig(t)
	err := commandExplore(cfg, []string{"explore", "bad-area"})
	if err == nil {
		t.Fatal("expected error from commandExplore when API fails")
	}
}

func TestCommandCatchRequiresArg(t *testing.T) {
	cfg := newTestConfig(t)
	if err := commandCatch(cfg, []string{"catch"}); err == nil {
		t.Fatal("expected error when no pokemon name is provided")
	}
}

func TestCommandCatchSuccess(t *testing.T) {
	// BaseExperience 0 ⇒ rand.IntN(maxCatchRate) >= 0 is always true ⇒ always caught.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"name":"weedle","base_experience":0}`))
	}))
	defer srv.Close()
	defer pokeapi.SetBaseURLForTest(srv.URL + "/")()

	cfg := newTestConfig(t)
	out := captureStdout(t, func() {
		if err := commandCatch(cfg, []string{"catch", "weedle"}); err != nil {
			t.Fatalf("commandCatch returned error: %v", err)
		}
	})
	if !strings.Contains(out, "was caught") {
		t.Errorf("expected catch success message, got: %q", out)
	}
	if _, ok := cfg.CaughtPokemonMap.Get("weedle"); !ok {
		t.Error("expected weedle to be added to caught map")
	}
}

func TestCommandCatchEscape(t *testing.T) {
	// BaseExperience == maxCatchRate ⇒ rand.IntN(maxCatchRate) is always < it ⇒ always escapes.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"name":"mewtwo","base_experience":700}`))
	}))
	defer srv.Close()
	defer pokeapi.SetBaseURLForTest(srv.URL + "/")()

	cfg := newTestConfig(t)
	out := captureStdout(t, func() {
		if err := commandCatch(cfg, []string{"catch", "mewtwo"}); err != nil {
			t.Fatalf("commandCatch returned error: %v", err)
		}
	})
	if !strings.Contains(out, "escaped") {
		t.Errorf("expected escape message, got: %q", out)
	}
	if _, ok := cfg.CaughtPokemonMap.Get("mewtwo"); ok {
		t.Error("expected mewtwo to NOT be added to caught map")
	}
}

func TestCommandCatchAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer srv.Close()
	defer pokeapi.SetBaseURLForTest(srv.URL + "/")()

	cfg := newTestConfig(t)
	err := commandCatch(cfg, []string{"catch", "missingno"})
	if err == nil {
		t.Fatal("expected error from commandCatch when API fails")
	}
}

func TestCommandInspectRequiresArg(t *testing.T) {
	cfg := newTestConfig(t)
	if err := commandInspect(cfg, []string{"inspect"}); err == nil {
		t.Fatal("expected error when no pokemon name is provided")
	}
}

func TestCommandInspectMissNotifiesUser(t *testing.T) {
	cfg := newTestConfig(t)
	out := captureStdout(t, func() {
		if err := commandInspect(cfg, []string{"inspect", "missingno"}); err != nil {
			t.Fatalf("commandInspect should not return error on miss, got: %v", err)
		}
	})
	if !strings.Contains(out, "have not caught") {
		t.Errorf("expected miss message, got: %q", out)
	}
}

func TestCommandInspectPrintsDetails(t *testing.T) {
	cfg := newTestConfig(t)
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

	out := captureStdout(t, func() {
		if err := commandInspect(cfg, []string{"inspect", "eevee"}); err != nil {
			t.Fatalf("commandInspect returned error: %v", err)
		}
	})
	for _, want := range []string{"eevee", "hp", "normal", "55"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got: %q", want, out)
		}
	}
}

func TestCommandPokedexEmpty(t *testing.T) {
	cfg := newTestConfig(t)
	if err := commandPokedex(cfg, nil); err == nil {
		t.Fatal("expected error when pokedex is empty")
	}
}

func TestCommandPokedexListsCaughtSorted(t *testing.T) {
	cfg := newTestConfig(t)
	cfg.CaughtPokemonMap.Add("zubat", pokeapi.PokemonDetails{Name: "zubat"})
	cfg.CaughtPokemonMap.Add("abra", pokeapi.PokemonDetails{Name: "abra"})
	cfg.CaughtPokemonMap.Add("magikarp", pokeapi.PokemonDetails{Name: "magikarp"})

	out := captureStdout(t, func() {
		if err := commandPokedex(cfg, nil); err != nil {
			t.Fatalf("commandPokedex returned error: %v", err)
		}
	})

	idxAbra := strings.Index(out, "abra")
	idxMagikarp := strings.Index(out, "magikarp")
	idxZubat := strings.Index(out, "zubat")
	if idxAbra < 0 || idxMagikarp < 0 || idxZubat < 0 {
		t.Fatalf("expected all three names in output, got: %q", out)
	}
	if !(idxAbra < idxMagikarp && idxMagikarp < idxZubat) {
		t.Errorf("expected alphabetical order, got positions abra=%d magikarp=%d zubat=%d", idxAbra, idxMagikarp, idxZubat)
	}
}
