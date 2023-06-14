package main

import (
	"flag"
)

var flagAddress string
var flagReportInterval int
var flagPollInterval int

func parseFlags() {
	flag.StringVar(&flagAddress, "a", "localhost:8080", "HTTP server endpoint address")
	flag.IntVar(&flagReportInterval, "r", 10, "input report interval")
	flag.IntVar(&flagPollInterval, "p", 2, "input poll interval")

	flag.Parse()
}
