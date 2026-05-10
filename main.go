package main

import (
	"os"

	"github.com/kelvinjrosado/pokedex/internal/repl"
)

func main() {
	// Run REPL
	repl.Run(os.Stdin, os.Stdout)
}
