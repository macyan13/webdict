package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/macyan13/webdict/backend/pkg/server"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Printf("[ERROR] failed with %+v", err)
		os.Exit(1)
	}
}

func run() error {
	var opts server.Opts
	if _, err := flags.Parse(&opts); err != nil {
		return err
	}

	s, err := server.InitServer(opts)
	if err != nil {
		return err
	}

	log.Printf("[INFO] starting server on port %s:%d", opts.WebdictURL, opts.Port)
	if err := s.Run(); err != nil {
		return err
	}

	return nil
}
