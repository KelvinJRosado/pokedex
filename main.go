package main

import (
	"log"
	"os"

	"github.com/kelvinjrosado/pokedex/internal/logger"
	"github.com/kelvinjrosado/pokedex/internal/repl"
)

func main() {
	lgr, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("Unable to start logger: %v", err)
	}
	defer lgr.Close()

	repl.Run(os.Stdin, os.Stdout, lgr)
}
