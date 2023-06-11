package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	pollInterval := 2 * time.Second
	reportInterval := 10 * time.Second

	pollTicker := time.NewTicker(pollInterval)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(reportInterval)
	defer reportTicker.Stop()

	metrics := Metrics{}

	client := http.Client{}

	for {
		select {
		case <-pollTicker.C:
			metrics = collectMetrics()

		case <-reportTicker.C:
			err := sendMetrics(&client, &metrics)
			if err != nil {
				log.Printf("Failed to send metrics: %v", err)
			}
		}
	}
}
