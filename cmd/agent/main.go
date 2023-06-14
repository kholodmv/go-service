package main

import (
	"github.com/go-resty/resty/v2"
	"log"
	"time"
)

func main() {
	parseFlags()

	pollTicker := time.NewTicker(flagPollInterval)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(flagReportInterval)
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
