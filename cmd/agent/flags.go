package main

import (
	"flag"
	"os"
	"strconv"
)

type AgentParams struct {
	flagAddress        string
	flagReportInterval int
	flagPollInterval   int
}

func useStartParams() AgentParams {
	f := AgentParams{}

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		f.flagAddress = envRunAddr
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		f.flagReportInterval, _ = strconv.Atoi(envReportInterval)
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		f.flagReportInterval, _ = strconv.Atoi(envPollInterval)
	}

	flag.StringVar(&f.flagAddress, "a", "localhost:8081", "HTTP server endpoint address")
	flag.IntVar(&f.flagReportInterval, "r", 10, "input report interval")
	flag.IntVar(&f.flagPollInterval, "p", 2, "input poll interval")

	flag.Parse()

	return f
}
