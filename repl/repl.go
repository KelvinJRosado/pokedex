package repl

import (
	"bufio"
	"fmt"
	"os"
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

// Runs the REPL
// Will read user input and reponds until session ended
func Run() {
	// Create scanner instance to read user input
	scanner := bufio.NewScanner(os.Stdin)

	// Print standard line
	fmt.Print("Pokedex > ")

	// Use an infinite loop to keep CLI open
	for scanner.Scan() {
		// Get user input
		input := scanner.Text()

		// Clean user input
		cleaned := cleanInput(input)

		// Check how much was entered
		if len(cleaned) > 0 {
			// Grab 1st word
			first := cleaned[0]

			// Print message
			fmt.Printf("Your command was: %v\n", first)
		}

		// Prepare next line
		fmt.Print("Pokedex > ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Invalid input: %s\n", err)
	}
}
