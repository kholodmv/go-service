// Package configs - agent_flags.go - agent configuration parameters.
package configs

import (
	"flag"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"
)

// AgentParams structure that contains flag variables.
type AgentParams struct {
	FlagAddress         string // flag HTTP server endpoint address
	FlagReportInterval  int    // flag report interval
	FlagPollInterval    int    // flag poll interval
	FlagKey             string // flag key for calculating SHA-256 hash
	FlagRateLimit       int    // flag rate limit
	FlagCryptoPublicKey string // flag crypto key
}

// ConfigAgent structure that contains variables for initial.
type ConfigAgent struct {
	Client          *resty.Client // resty client instance
	AgentURL        string        // HTTP server endpoint address
	ReportInterval  int           // report interval
	PollInterval    int           // poll interval
	Key             string        // key for calculating SHA-256 hash
	RateLimit       int           // rate limit
	CryptoPublicKey string        // crypto key
}

// InitConfigAgent - agent configuration initialization function.
func InitConfigAgent() ConfigAgent {
	f := AgentParams{}
	useAgentStartParams(&f)

	return ConfigAgent{
		Client:          resty.New(),
		AgentURL:        "http://" + f.FlagAddress + "/update/",
		ReportInterval:  f.FlagReportInterval,
		PollInterval:    f.FlagPollInterval,
		Key:             f.FlagKey,
		RateLimit:       f.FlagRateLimit,
		CryptoPublicKey: f.FlagCryptoPublicKey,
	}
}

// useAgentStartParams - assigning configuration environment variables.
func useAgentStartParams(f *AgentParams) {
	flag.StringVar(&f.FlagAddress, "a", "localhost:8080", "HTTP server endpoint address")
	flag.IntVar(&f.FlagReportInterval, "r", 10, "input report interval")
	flag.IntVar(&f.FlagPollInterval, "p", 2, "input poll interval")
	flag.StringVar(&f.FlagKey, "k", "", "KEY for calculating SHA-256 hash")
	flag.IntVar(&f.FlagRateLimit, "l", 3, "rate limit")
	flag.StringVar(&f.FlagCryptoPublicKey, "crypto-key", "", "path to RSA public key file in PEM format")
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
	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		f.FlagRateLimit, _ = strconv.Atoi(envRateLimit)
	}
	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		f.FlagCryptoPublicKey = envCryptoKey
	}
}
