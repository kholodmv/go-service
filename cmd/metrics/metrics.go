// Package metrics - metrics.go - agent work with metrics.
package metrics

import (
	"bytes"
	"compress/gzip"
	"context"
	cr "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/kholodmv/go-service/internal/middleware/logger"
	"github.com/kholodmv/go-service/proto"
	"google.golang.org/grpc/metadata"
	"math/rand"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/kholodmv/go-service/internal/configs"
	"github.com/kholodmv/go-service/internal/models"
)

// Metrics struct includes a collection of metrics and
// a cumulative value for the PollCount metric with count type.
type Metrics struct {
	data      []models.Metrics
	pollCount int64
}

// log - initializing the logging variable.
var log = logger.Initialize()

// Compress - compress data.
// If successful, the compressed data is returned, else - error.
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

// ReportAgent function that triggers the collection and sending of metrics.
func (m *Metrics) ReportAgent(ctx context.Context, connectionsClosed chan struct{}, c configs.ConfigAgent, pk *rsa.PublicKey, metricClient proto.MetricsClient) {
	metricCh := make(chan models.Metrics)
	timeR := 0
	for {
		if timeR >= c.ReportInterval {
			timeR = 0
			go m.SendMetrics(c.Client, c.AgentURL, c.Key, metricCh, c.RateLimit, pk)
			go m.sendGRPC(metricCh, metricClient)
		}
		go m.CollectMetrics(metricCh)

		time.Sleep(time.Duration(c.PollInterval) * time.Second)
		timeR += c.PollInterval
	}
}

// CollectMetrics function collecting metrics of different types.
func (m *Metrics) CollectMetrics(ch chan<- models.Metrics) {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)

	var v float64

	v = float64(memStats.Alloc)
	ch <- models.Metrics{ID: "Alloc", MType: "gauge", Value: &v}
	v = float64(memStats.BuckHashSys)
	ch <- models.Metrics{ID: "BuckHashSys", MType: "gauge", Value: &v}
	v = float64(memStats.Frees)
	ch <- models.Metrics{ID: "Frees", MType: "gauge", Value: &v}
	v = memStats.GCCPUFraction
	ch <- models.Metrics{ID: "GCCPUFraction", MType: "gauge", Value: &v}
	v = float64(memStats.GCSys)
	ch <- models.Metrics{ID: "GCSys", MType: "gauge", Value: &v}
	v = float64(memStats.HeapAlloc)
	ch <- models.Metrics{ID: "HeapAlloc", MType: "gauge", Value: &v}
	v = float64(memStats.HeapIdle)
	ch <- models.Metrics{ID: "HeapIdle", MType: "gauge", Value: &v}
	v = float64(memStats.HeapInuse)
	ch <- models.Metrics{ID: "HeapInuse", MType: "gauge", Value: &v}
	v = float64(memStats.HeapObjects)
	ch <- models.Metrics{ID: "HeapObjects", MType: "gauge", Value: &v}
	v = float64(memStats.HeapReleased)
	ch <- models.Metrics{ID: "HeapReleased", MType: "gauge", Value: &v}
	v = float64(memStats.HeapSys)
	ch <- models.Metrics{ID: "HeapSys", MType: "gauge", Value: &v}
	v = float64(memStats.LastGC)
	ch <- models.Metrics{ID: "LastGC", MType: "gauge", Value: &v}
	v = float64(memStats.Lookups)
	ch <- models.Metrics{ID: "Lookups", MType: "gauge", Value: &v}
	v = float64(memStats.MCacheInuse)
	ch <- models.Metrics{ID: "MCacheInuse", MType: "gauge", Value: &v}
	v = float64(memStats.MCacheSys)
	ch <- models.Metrics{ID: "MCacheSys", MType: "gauge", Value: &v}
	v = float64(memStats.MSpanInuse)
	ch <- models.Metrics{ID: "MSpanInuse", MType: "gauge", Value: &v}
	v = float64(memStats.MSpanSys)
	ch <- models.Metrics{ID: "MSpanSys", MType: "gauge", Value: &v}
	v = float64(memStats.Mallocs)
	ch <- models.Metrics{ID: "Mallocs", MType: "gauge", Value: &v}
	v = float64(memStats.NextGC)
	ch <- models.Metrics{ID: "NextGC", MType: "gauge", Value: &v}
	v = float64(memStats.NumForcedGC)
	ch <- models.Metrics{ID: "NumForcedGC", MType: "gauge", Value: &v}
	v = float64(memStats.NumGC)
	ch <- models.Metrics{ID: "NumGC", MType: "gauge", Value: &v}
	v = float64(memStats.OtherSys)
	ch <- models.Metrics{ID: "OtherSys", MType: "gauge", Value: &v}
	v = float64(memStats.PauseTotalNs)
	ch <- models.Metrics{ID: "PauseTotalNs", MType: "gauge", Value: &v}
	v = float64(memStats.StackInuse)
	ch <- models.Metrics{ID: "StackInuse", MType: "gauge", Value: &v}
	v = float64(memStats.StackSys)
	ch <- models.Metrics{ID: "StackSys", MType: "gauge", Value: &v}
	v = float64(memStats.Sys)
	ch <- models.Metrics{ID: "Sys", MType: "gauge", Value: &v}
	v = float64(memStats.TotalAlloc)
	ch <- models.Metrics{ID: "TotalAlloc", MType: "gauge", Value: &v}

	vv, _ := mem.VirtualMemory()
	v = float64(vv.Total)
	ch <- models.Metrics{ID: "TotalMemory", MType: "gauge", Value: &v}
	v = float64(vv.Free)
	ch <- models.Metrics{ID: "FreeMemory", MType: "gauge", Value: &v}

	cpuInfo, _ := cpu.Info()
	numCPUs := float64(len(cpuInfo))
	cpuUtil, _ := cpu.Percent(time.Second, false)
	v = cpuUtil[0] / numCPUs
	ch <- models.Metrics{ID: "CPUutilization1", MType: "gauge", Value: &v}

	r := rand.Float64()
	ch <- models.Metrics{ID: "RandomValue", MType: "gauge", Value: &r}

	i := m.pollCount
	ch <- models.Metrics{ID: "PollCount", MType: "counter", Delta: &i}

	m.pollCount += 1
}

// SendMetrics function sending metrics by URL.
func (m *Metrics) SendMetrics(client *resty.Client, agentURL string, key string, metricCh <-chan models.Metrics, rateLimit int, pk *rsa.PublicKey) error {
	for metric := range metricCh {
		url := agentURL

		metricsJSON, err := json.Marshal(metric)
		if err != nil {
			fmt.Printf("Error metrics JSON: %s\n", err)
		}
		hashSHA256 := fmt.Sprintf("%x", sha256.Sum256([]byte(metricsJSON)))
		metricsJSON, err = Compress(metricsJSON)
		if err != nil {
			fmt.Printf("Error compress JSON: %s\n", err)
		}

		encryptedBytes, err := rsa.EncryptOAEP(
			sha256.New(),
			cr.Reader,
			pk,
			metricsJSON,
			nil)

		if err != nil {
			return err
		}
		reader := bytes.NewReader(encryptedBytes)

		ip, err := GetLocalIP()
		if err != nil {
			return err
		}

		var resp *resty.Response
		if key != "" {
			resp, err = client.R().
				SetBody(reader).
				SetHeader("Content-Type", "application/json").
				SetHeader("Accept", "application/json").
				SetHeader("Content-Encoding", "gzip").
				SetHeader("HashSHA256", hashSHA256).
				SetHeader("X-Real-IP", ip).
				Post(url)
		} else {
			resp, err = client.R().
				SetBody(reader).
				SetHeader("Content-Type", "application/json").
				SetHeader("Accept", "application/json").
				SetHeader("Content-Encoding", "gzip").
				SetHeader("X-Real-IP", ip).
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

		time.Sleep(time.Second)

		ticker := time.NewTicker(time.Duration(rateLimit) * time.Second)
		<-ticker.C
		ticker.Stop()
	}

	return nil
}

func (m *Metrics) sendGRPC(metricCh <-chan models.Metrics, metricClient proto.MetricsClient) error {
	for metric := range metricCh {
		pb, err := metric.ToProto()
		if err != nil {
			return err
		}

		ip, err := GetLocalIP()
		if err != nil {
			return err
		}
		metadata.AppendToOutgoingContext(context.Background(), "X-Real-IP", ip)

		_, err = metricClient.UpdateList(context.Background(), &proto.UpdateListRequest{Metric: pb})
		if err != nil {
			return errors.Wrap(err, "unable to make grpc call")
		}
	}
	return nil
}

func GetLocalIP() (string, error) {
	ips, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for i := range ips {
		if ips[i].String() != "127.0.0.1/8" {
			ip, _, err := net.ParseCIDR(ips[i].String())
			if err != nil {
				return "", err
			}

			return ip.String(), nil
		}
	}

	return "", errors.New("no IP available")
}
