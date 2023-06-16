package main

import (
	"github.com/go-resty/resty/v2"
	"log"
	"time"
)

func main() {
	flags := useStartParams()

	pollInterval := time.Duration(flags.flagPollInterval) * time.Second
	pollTicker := time.NewTicker(pollInterval)
	defer pollTicker.Stop()

	reportInterval := time.Duration(flags.flagReportInterval) * time.Second
	reportTicker := time.NewTicker(reportInterval)
	defer reportTicker.Stop()

	metrics := Metrics{}

	client := resty.New()

	agentURL := "http://" + flags.flagAddress + "/update"
	for {
		select {
		case <-pollTicker.C:
			metrics = collectMetrics()

		case <-reportTicker.C:
			err := sendMetrics(client, &metrics, agentURL)
			if err != nil {
				log.Printf("Failed to send metrics: %v", err)
			}
		}
	}
}
