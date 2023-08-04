package configs

import (
	"flag"
	"github.com/go-resty/resty/v2"
	"os"
	"strconv"
)

type AgentParams struct {
	FlagAddress        string
	FlagReportInterval int
	FlagPollInterval   int
	FlagKey            string
}

type ConfigAgent struct {
	Client         *resty.Client
	AgentURL       string
	ReportInterval int
	PollInterval   int
	Key            string
}

func InitConfigAgent() ConfigAgent {
	f := AgentParams{}
	useAgentStartParams(&f)

	return ConfigAgent{
		Client:         resty.New(),
		AgentURL:       "http://" + f.FlagAddress + "/update/",
		ReportInterval: f.FlagReportInterval,
		PollInterval:   f.FlagPollInterval,
		Key:            f.FlagKey,
	}
}

func useAgentStartParams(f *AgentParams) {
	flag.StringVar(&f.FlagAddress, "a", "localhost:8080", "HTTP server endpoint address")
	flag.IntVar(&f.FlagReportInterval, "r", 10, "input report interval")
	flag.IntVar(&f.FlagPollInterval, "p", 2, "input poll interval")
	flag.StringVar(&f.FlagKey, "k", "kkk", "KEY for calculating SHA-256 hash")

	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		f.FlagAddress = envRunAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		f.FlagReportInterval, _ = strconv.Atoi(envReportInterval)
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		f.FlagPollInterval, _ = strconv.Atoi(envPollInterval)
	}
	if envKey := os.Getenv("KEY"); envKey != "" {
		f.FlagKey = envKey
	}
}
