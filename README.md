# Pokedex

A command-line Pokedex written in Go. Launches an interactive REPL that talks to the [PokeAPI](https://pokeapi.co/) so you can browse the Pokemon world, catch Pokemon, and inspect what you've collected.

Initially built as part of the [Boot.dev](https://www.boot.dev/) backend curriculum, but with additional enhancements.

## Requirements

- Go 1.26 or newer

## Run

```sh
go run .
```

You'll be dropped into the `Pokedex >` prompt. Type `help` to see the available commands.

## Build

```sh
go build -o pokedex
./pokedex
```

## Commands

| Command | Description |
| --- | --- |
| `help` | Show the command list |
| `map` | Show the next 20 locations |
| `mapb` | Show the previous 20 locations |
| `explore <location>` | List Pokemon that can be encountered in a location |
| `catch <pokemon>` | Try to catch a Pokemon (success depends on its base experience) |
| `inspect <pokemon>` | Show stats and types for a Pokemon you've caught |
| `pokedex` | List every Pokemon you've caught this session |
| `exit` | Quit |

## Project layout

```
.
├── main.go                  // entry point — boots the REPL
└── internal/
    ├── repl/                // REPL loop, command registry, input parsing
    ├── pokeapi/             // PokeAPI client + response models
    └── pokecache/           // in-memory cache with periodic reaping
```

API responses are cached in-memory with a background reaper, so paging back and forth through the map or re-exploring an area doesn't re-hit the network.

## Tests

```sh
go test ./...
```
