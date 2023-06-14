package main

import (
	"github.com/go-resty/resty/v2"
	"log"
	"time"
)

func main() {
	parseFlags()

	pollInterval := time.Duration(flagPollInterval) * time.Second
	pollTicker := time.NewTicker(pollInterval)
	defer pollTicker.Stop()

	reportInterval := time.Duration(flagReportInterval) * time.Second
	reportTicker := time.NewTicker(reportInterval)
	defer reportTicker.Stop()

	metrics := Metrics{}

	client := resty.New()

	for {
		select {
		case <-pollTicker.C:
			metrics = collectMetrics()

		case <-reportTicker.C:
			err := sendMetrics(client, &metrics, flagAddress)
			if err != nil {
				log.Printf("Failed to send metrics: %v", err)
			}
		}
	}
}
