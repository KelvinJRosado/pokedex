package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
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

		// Grab 1st word
		first := cleaned[0]

		// Print message
		fmt.Printf("Your command was: %v\n", first)

		// Prepare next line
		fmt.Print("Pokedex > ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Invalid input: %s\n", err)
	}
}
