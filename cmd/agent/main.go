package main

import (
	"github.com/go-resty/resty/v2"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	parseFlags()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagAddress = envRunAddr
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		flagReportInterval, _ = strconv.Atoi(envReportInterval)
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		flagReportInterval, _ = strconv.Atoi(envPollInterval)
	}

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
