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
}

type ConfigAgent struct {
	Client         *resty.Client
	AgentUrl       string
	ReportInterval int
	PollInterval   int
}

func InitConfigAgent() ConfigAgent {
	f := AgentParams{}
	useAgentStartParams(&f)

	return ConfigAgent{
		Client:         resty.New(),
		AgentUrl:       "http://" + f.FlagAddress + "/update",
		ReportInterval: f.FlagReportInterval,
		PollInterval:   f.FlagPollInterval,
	}
}

func useAgentStartParams(f *AgentParams) {
	flag.StringVar(&f.FlagAddress, "a", "localhost:8080", "HTTP server endpoint address")
	flag.IntVar(&f.FlagReportInterval, "r", 10, "input report interval")
	flag.IntVar(&f.FlagPollInterval, "p", 2, "input poll interval")

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
}
