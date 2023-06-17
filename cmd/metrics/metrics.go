package metrics

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
)

type gauge float64
type counter int64

type Metrics struct {
	runtimeMetrics map[string]gauge
	randomValue    gauge
	pollCount      counter
}

func CollectMetrics() Metrics {
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
	r := rand.New(rand.NewSource(99))
	metrics.randomValue = gauge(r.Int63())

	return metrics
}

func (m *Metrics) SendMetrics(client *resty.Client, agentURL string) error {
	for k, v := range m.runtimeMetrics {
		url := agentURL + "/gauge/" + k + "/" + strconv.FormatFloat(float64(v), 'f', 1, 64)

		resp, err := client.R().
			SetHeader("Content-Type", "text/plain").
			Post(url)

		if err != nil {
			return fmt.Errorf("HTTP POST request failed: %v", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("unexpected HTTP response status")
		}

		log.Println(url)
		log.Println(resp)
		url = ""
	}

	url := agentURL + "/counter/someMetric/" + strconv.FormatInt(int64(m.pollCount), 10)
	resp, err := client.R().
		SetHeader("Content-Type", "text/plain").
		Post(url)
	if err != nil {
		return fmt.Errorf("HTTP POST request failed: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected HTTP response status")
	}

	log.Println(url)
	log.Println(resp)

	return nil
}
