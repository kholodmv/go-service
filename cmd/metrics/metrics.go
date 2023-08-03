package metrics

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kholodmv/go-service/internal/configs"
	"github.com/kholodmv/go-service/internal/logger"
	"github.com/kholodmv/go-service/internal/models"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type Metrics struct {
	data      []models.Metrics
	pollCount int64
}

var log = logger.Initialize()

func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("failed init compress writer: %v", err)
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	return b.Bytes(), nil
}

func (m *Metrics) ReportAgent(c configs.ConfigAgent) {
	timeR := 0
	for {
		if timeR >= c.ReportInterval {
			timeR = 0
			err := m.SendMetrics(c.Client, c.AgentURL, c.Key)
			if err != nil {
				log.Infow("Failed to send metrics: %v", err)
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

	r := rand.Float64()
	m.data = append(m.data, models.Metrics{ID: "RandomValue", MType: "gauge", Value: &r})

	i := m.pollCount
	m.data = append(m.data, models.Metrics{ID: "PollCount", MType: "counter", Delta: &i})

	m.pollCount += 1
}

func (m *Metrics) SendMetrics(client *resty.Client, agentURL string, key string) error {
	for _, metrics := range m.data {
		url := agentURL

		metricsJSON, err := json.Marshal(metrics)
		if err != nil {
			fmt.Printf("Error metrics JSON: %s\n", err)
		}
		metricsJSON, err = Compress(metricsJSON)
		if err != nil {
			fmt.Printf("Error compress JSON: %s\n", err)
		}

		var resp *resty.Response
		if key != "" {
			hashedKey := calculateSHA256(key)
			resp, err = client.R().
				SetBody(metricsJSON).
				SetHeader("Accept", "application/json").
				SetHeader("Accept-Encoding", "gzip").
				SetHeader("Content-Type", "application/json").
				SetHeader("HashSHA256", hashedKey).
				Post(url)
		} else {
			resp, err = client.R().
				SetBody(metricsJSON).
				SetHeader("Accept", "application/json").
				SetHeader("Accept-Encoding", "gzip").
				SetHeader("Content-Type", "application/json").
				SetHeader("Content-Type", "application/json").
				Post(url)
		}

		if err != nil {
			return fmt.Errorf("HTTP POST request failed: %v", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("unexpected HTTP response status")
		}

		log.Info(url)
		log.Info(resp)
		url = ""
	}

	return nil
}

func calculateSHA256(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}
