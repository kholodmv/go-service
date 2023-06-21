package metrics

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kholodmv/go-service/internal/configs"
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

func (m *Metrics) ReportAgent(c configs.ConfigAgent) {
	timeR := 0
	for {
		if timeR >= c.ReportInterval {
			timeR = 0
			err := m.SendMetrics(c.Client, c.AgentURL)
			if err != nil {
				log.Printf("Failed to send metrics: %v", err)
			}
		}
		m.CollectMetrics()

		time.Sleep(time.Duration(c.PollInterval) * time.Second)
		timeR += c.PollInterval
	}
}

func (m *Metrics) CollectMetrics() {
	m.runtimeMetrics = make(map[string]gauge)

	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)

	m.runtimeMetrics["Alloc"] = gauge(memStats.Alloc)
	m.runtimeMetrics["BuckHashSys"] = gauge(memStats.BuckHashSys)
	m.runtimeMetrics["Frees"] = gauge(memStats.Frees)
	m.runtimeMetrics["GCCPUFraction"] = gauge(memStats.GCCPUFraction)
	m.runtimeMetrics["GCSys"] = gauge(memStats.GCSys)
	m.runtimeMetrics["HeapAlloc"] = gauge(memStats.HeapAlloc)
	m.runtimeMetrics["HeapIdle"] = gauge(memStats.HeapIdle)
	m.runtimeMetrics["HeapInuse"] = gauge(memStats.HeapInuse)
	m.runtimeMetrics["HeapObjects"] = gauge(memStats.HeapObjects)
	m.runtimeMetrics["HeapReleased"] = gauge(memStats.HeapReleased)
	m.runtimeMetrics["HeapSys"] = gauge(memStats.HeapSys)
	m.runtimeMetrics["LastGC"] = gauge(memStats.LastGC)
	m.runtimeMetrics["Lookups"] = gauge(memStats.Lookups)
	m.runtimeMetrics["MCacheInuse"] = gauge(memStats.MCacheInuse)
	m.runtimeMetrics["MCacheSys"] = gauge(memStats.MCacheSys)
	m.runtimeMetrics["MSpanInuse"] = gauge(memStats.MSpanInuse)
	m.runtimeMetrics["MSpanSys"] = gauge(memStats.MSpanSys)
	m.runtimeMetrics["Mallocs"] = gauge(memStats.Mallocs)
	m.runtimeMetrics["NextGC"] = gauge(memStats.NextGC)
	m.runtimeMetrics["NumForcedGC"] = gauge(memStats.NumForcedGC)
	m.runtimeMetrics["NumGC"] = gauge(memStats.NumGC)
	m.runtimeMetrics["OtherSys"] = gauge(memStats.OtherSys)
	m.runtimeMetrics["PauseTotalNs"] = gauge(memStats.PauseTotalNs)
	m.runtimeMetrics["StackInuse"] = gauge(memStats.StackInuse)
	m.runtimeMetrics["StackSys"] = gauge(memStats.StackSys)
	m.runtimeMetrics["Sys"] = gauge(memStats.Sys)
	m.runtimeMetrics["TotalAlloc"] = gauge(memStats.TotalAlloc)

	m.pollCount += 1

	r := rand.New(rand.NewSource(99))
	m.randomValue = gauge(r.Int63())
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
