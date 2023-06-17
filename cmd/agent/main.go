package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/kholodmv/go-service/cmd/metrics"
	"github.com/kholodmv/go-service/internal/configs"
	"log"
	"time"
)

func main() {
	flags := configs.UseAgentStartParams()

	pollInterval := time.Duration(flags.FlagPollInterval) * time.Second
	pollTicker := time.NewTicker(pollInterval)
	defer pollTicker.Stop()

	reportInterval := time.Duration(flags.FlagReportInterval) * time.Second
	reportTicker := time.NewTicker(reportInterval)
	defer reportTicker.Stop()

	m := metrics.Metrics{}

	client := resty.New()

	agentURL := "http://" + flags.FlagAddress + "/update"
	for {
		select {
		case <-pollTicker.C:
			m = metrics.CollectMetrics()

		case <-reportTicker.C:
			err := m.SendMetrics(client, agentURL)
			if err != nil {
				log.Printf("Failed to send metrics: %v", err)
			}
		}
	}
}
