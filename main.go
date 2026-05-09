package main

import (
	"os"

	"github.com/kelvinjrosado/pokedex/internal/repl"
)

func main() {
	repl.Run(os.Stdin, os.Stdout)
}
