package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

type gauge float64
type counter int64

type Metrics struct {
	runtimeMetrics map[string]gauge
	randomValue    gauge
	pollCount      counter
}

func collectMetrics() Metrics {
	metrics := Metrics{
		runtimeMetrics: make(map[string]gauge),
	}
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)

	metrics.runtimeMetrics["Alloc"] = gauge(memStats.Alloc)
	metrics.runtimeMetrics["BuckHashSys"] = gauge(memStats.BuckHashSys)
	metrics.runtimeMetrics["Frees"] = gauge(memStats.Frees)
	metrics.runtimeMetrics["GCCPUFraction"] = gauge(memStats.GCCPUFraction)
	metrics.runtimeMetrics["GCSys"] = gauge(memStats.GCSys)
	metrics.runtimeMetrics["HeapAlloc"] = gauge(memStats.HeapAlloc)
	metrics.runtimeMetrics["HeapIdle"] = gauge(memStats.HeapIdle)
	metrics.runtimeMetrics["HeapInuse"] = gauge(memStats.HeapInuse)
	metrics.runtimeMetrics["HeapObjects"] = gauge(memStats.HeapObjects)
	metrics.runtimeMetrics["HeapReleased"] = gauge(memStats.HeapReleased)
	metrics.runtimeMetrics["HeapSys"] = gauge(memStats.HeapSys)
	metrics.runtimeMetrics["LastGC"] = gauge(memStats.LastGC)
	metrics.runtimeMetrics["Lookups"] = gauge(memStats.Lookups)
	metrics.runtimeMetrics["MCacheInuse"] = gauge(memStats.MCacheInuse)
	metrics.runtimeMetrics["MCacheSys"] = gauge(memStats.MCacheSys)
	metrics.runtimeMetrics["MSpanInuse"] = gauge(memStats.MSpanInuse)
	metrics.runtimeMetrics["MSpanSys"] = gauge(memStats.MSpanSys)
	metrics.runtimeMetrics["Mallocs"] = gauge(memStats.Mallocs)
	metrics.runtimeMetrics["NextGC"] = gauge(memStats.NextGC)
	metrics.runtimeMetrics["NumForcedGC"] = gauge(memStats.NumForcedGC)
	metrics.runtimeMetrics["NumGC"] = gauge(memStats.NumGC)
	metrics.runtimeMetrics["OtherSys"] = gauge(memStats.OtherSys)
	metrics.runtimeMetrics["PauseTotalNs"] = gauge(memStats.PauseTotalNs)
	metrics.runtimeMetrics["StackInuse"] = gauge(memStats.StackInuse)
	metrics.runtimeMetrics["StackSys"] = gauge(memStats.StackSys)
	metrics.runtimeMetrics["Sys"] = gauge(memStats.Sys)
	metrics.runtimeMetrics["TotalAlloc"] = gauge(memStats.TotalAlloc)

	metrics.pollCount += 1

	rand.Seed(time.Now().UnixNano())
	metrics.randomValue = gauge(rand.Int63())

	return metrics
}

func sendMetrics(client *http.Client, metrics *Metrics) error {
	domain := "http://localhost:8080"

	for k, v := range metrics.runtimeMetrics {
		url := domain + "/update/gauge/" + k + "/" + strconv.FormatFloat(float64(v), 'f', 1, 64)
		resp, err := client.Post(url, "text/plain", nil)
		if err != nil {
			return fmt.Errorf("HTTP POST request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected HTTP response status: %s", resp.Status)
		}

		log.Println(resp)
		url = ""
	}
	return nil
}
