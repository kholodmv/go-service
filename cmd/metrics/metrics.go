package metrics

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kholodmv/go-service/internal/configs"
	"github.com/kholodmv/go-service/internal/models"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type Metrics struct {
	data      []models.Metrics
	pollCount int64
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
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)

	m.data = []models.Metrics{}

	var v float64

	v = float64(memStats.Alloc)
	m.data = append(m.data, models.Metrics{ID: "Alloc", MType: "gauge", Value: &v})
	v = float64(memStats.BuckHashSys)
	m.data = append(m.data, models.Metrics{ID: "BuckHashSys", MType: "gauge", Value: &v})
	v = float64(memStats.Frees)
	m.data = append(m.data, models.Metrics{ID: "Frees", MType: "gauge", Value: &v})
	v = memStats.GCCPUFraction
	m.data = append(m.data, models.Metrics{ID: "GCCPUFraction", MType: "gauge", Value: &v})
	v = float64(memStats.GCSys)
	m.data = append(m.data, models.Metrics{ID: "GCSys", MType: "gauge", Value: &v})
	v = float64(memStats.HeapAlloc)
	m.data = append(m.data, models.Metrics{ID: "HeapAlloc", MType: "gauge", Value: &v})
	v = float64(memStats.HeapIdle)
	m.data = append(m.data, models.Metrics{ID: "HeapIdle", MType: "gauge", Value: &v})
	v = float64(memStats.HeapInuse)
	m.data = append(m.data, models.Metrics{ID: "HeapInuse", MType: "gauge", Value: &v})
	v = float64(memStats.HeapObjects)
	m.data = append(m.data, models.Metrics{ID: "HeapObjects", MType: "gauge", Value: &v})
	v = float64(memStats.HeapReleased)
	m.data = append(m.data, models.Metrics{ID: "HeapReleased", MType: "gauge", Value: &v})
	v = float64(memStats.HeapSys)
	m.data = append(m.data, models.Metrics{ID: "HeapSys", MType: "gauge", Value: &v})
	v = float64(memStats.LastGC)
	m.data = append(m.data, models.Metrics{ID: "LastGC", MType: "gauge", Value: &v})
	v = float64(memStats.Lookups)
	m.data = append(m.data, models.Metrics{ID: "Lookups", MType: "gauge", Value: &v})
	v = float64(memStats.MCacheInuse)
	m.data = append(m.data, models.Metrics{ID: "MCacheInuse", MType: "gauge", Value: &v})
	v = float64(memStats.MCacheSys)
	m.data = append(m.data, models.Metrics{ID: "MCacheSys", MType: "gauge", Value: &v})
	v = float64(memStats.MSpanInuse)
	m.data = append(m.data, models.Metrics{ID: "MSpanInuse", MType: "gauge", Value: &v})
	v = float64(memStats.MSpanSys)
	m.data = append(m.data, models.Metrics{ID: "MSpanSys", MType: "gauge", Value: &v})
	v = float64(memStats.Mallocs)
	m.data = append(m.data, models.Metrics{ID: "Mallocs", MType: "gauge", Value: &v})
	v = float64(memStats.NextGC)
	m.data = append(m.data, models.Metrics{ID: "NextGC", MType: "gauge", Value: &v})
	v = float64(memStats.NumForcedGC)
	m.data = append(m.data, models.Metrics{ID: "NumForcedGC", MType: "gauge", Value: &v})
	v = float64(memStats.NumGC)
	m.data = append(m.data, models.Metrics{ID: "NumGC", MType: "gauge", Value: &v})
	v = float64(memStats.OtherSys)
	m.data = append(m.data, models.Metrics{ID: "OtherSys", MType: "gauge", Value: &v})
	v = float64(memStats.PauseTotalNs)
	m.data = append(m.data, models.Metrics{ID: "PauseTotalNs", MType: "gauge", Value: &v})
	v = float64(memStats.StackInuse)
	m.data = append(m.data, models.Metrics{ID: "StackInuse", MType: "gauge", Value: &v})
	v = float64(memStats.StackSys)
	m.data = append(m.data, models.Metrics{ID: "StackSys", MType: "gauge", Value: &v})
	v = float64(memStats.Sys)
	m.data = append(m.data, models.Metrics{ID: "Sys", MType: "gauge", Value: &v})
	v = float64(memStats.TotalAlloc)
	m.data = append(m.data, models.Metrics{ID: "TotalAlloc", MType: "gauge", Value: &v})

	r := rand.New(rand.NewSource(99))
	m.data = append(m.data, models.Metrics{ID: "RandomValue", MType: "gauge", Value: &r})

	m.data = append(m.data, models.Metrics{ID: "PollCount", MType: "counter", Delta: &m.pollCount})

	m.pollCount += 1
}

func (m *Metrics) SendMetrics(client *resty.Client, agentURL string) error {
	for _, metrics := range m.data {
		url := agentURL

		metricsJSON, err := json.Marshal(metrics)
		if err != nil {
			fmt.Printf("Error metrics JSON: %s\n", err)
		}

		resp, err := client.R().
			SetBody(metricsJSON).
			SetHeader("Content-Type", "application/json").
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

	return nil
}
