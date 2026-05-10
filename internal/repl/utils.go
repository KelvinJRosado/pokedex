package repl

import (
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, []string) error
}

// Split the user's input into words based on whitespace.
// It should also lowercase the input and trim any leading or trailing whitespace.
func cleanInput(text string) []string {
	var out []string

	// First we make everything lowercase
	lower := strings.ToLower(text)

	// Lastly we split on whitespace, removing leading/trailing whitespace
	out = strings.Fields(lower)

	return out
}

// getAllCommands returns the command registry as a fresh slice on every call.
// A slice is used instead of a map so there is no package-level mutable state that
// could be accidentally modified. The slice is rebuilt on each call, making the
// registry effectively immutable — callers get their own copy and cannot mutate
// a shared reference. This also avoids the initialization cycle that occurs when
// a package-level variable (map or slice) references command functions whose
// bodies in turn read from that same variable.
//
// The tradeoff is that this allocates a new slice on every invocation, whereas a
// package-level map or slice would allocate once and reuse. However, with only a
// handful of commands in the registry, the allocation is trivial and well worth
// the guarantee of no shared mutable state.
func getAllCommands() []cliCommand {
	return []cliCommand{
		{
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		{
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		{
			name:        "map",
			description: "Displays the name of the 20 next locations in the Pokemon world",
			callback:    commandMap,
		},
		{
			name:        "mapb",
			description: "Displays the name of the 20 previous locations in the Pokemon world",
			callback:    commandMapb,
		},
		{
			name:        "explore",
			description: "Displays the name of Pokemon that can be encountered in the specified area",
			callback:    commandExplore,
		},
		{
			name:        "catch",
			description: "Attempt to catch the specified Pokemon",
			callback:    commandCatch,
		},
		{
			name:        "inspect",
			description: "Displays information about the specified Pokemon if caught",
			callback:    commandInspect,
		},
		{
			name:        "pokedex",
			description: "Displays all caught Pokemon",
			callback:    commandPokedex,
		},
	}
}

func getCommand(name string) (cliCommand, bool) {
	for _, cmd := range getAllCommands() {
		if name == cmd.name {
			return cmd, true
		}
	}
	return cliCommand{}, false
}
