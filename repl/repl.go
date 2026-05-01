package repl

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

// Runs the REPL
// Will read user input and reponds until session ended
func Run() {
	// Start the command registry
	initRegistry()

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

			// Lookup registry entry
			command, ok := commandRegistry[first]
			if !ok {
				fmt.Println("Unknown command")
				fmt.Print("Pokedex > ")
				continue
			}

			// Do command
			err := command.callback()
			if err != nil {
				// Check for clean exit
				if errors.Is(err, CleanExit) {
					os.Exit(0)
				}

				fmt.Printf("Error: %v", err.Error())
			}

		}

		// Prepare next line
		fmt.Print("Pokedex > ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Invalid input: %s\n", err)
	}
}
