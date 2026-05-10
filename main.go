package main

import (
	"log"
	"os"

	"github.com/kelvinjrosado/pokedex/internal/logger"
	"github.com/kelvinjrosado/pokedex/internal/repl"
)

func main() {
	// Create logger
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatal("Unable to start logger")
	}
	defer logger.Close()

	// Run REPL
	repl.Run(os.Stdin, os.Stdout, *logger)
}
