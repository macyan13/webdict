package main

import (
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/server"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "this is the startup error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	s := server.InitServer()
	err := s.Run()

	if err != nil {
		return err
	}

	return nil
}
