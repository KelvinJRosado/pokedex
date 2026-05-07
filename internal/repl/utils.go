package repl

import (
	"strings"
)

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
