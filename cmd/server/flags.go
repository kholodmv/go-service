package main

import (
	"flag"
	"os"
)

func useStartParams() string {
	var flagRunAddr string

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	return flagRunAddr
}
