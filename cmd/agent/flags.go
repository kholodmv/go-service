package main

import (
	"flag"
	"time"
)

var flagAddress string
var flagReportInterval time.Duration
var flagPollInterval time.Duration

func parseFlags() {
	flag.StringVar(&flagAddress, "a", "localhost:8080", "HTTP server endpoint address")
	flag.DurationVar(&flagReportInterval, "r", 10*time.Second, "input report interval")
	flag.DurationVar(&flagPollInterval, "p", 2*time.Second, "input poll interval")

	flag.Parse()
}
