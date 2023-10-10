// Package store - memStore.go - contains the logic for working with the memory.
package store

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/kholodmv/go-service/cmd/metrics"
	"github.com/kholodmv/go-service/internal/models"
)

// memoryStorage struct contains a collections of metrics and mutex.
type memoryStorage struct {
	mu             sync.Mutex
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
}

// NewMemoryStorage function which initialize collections of metrics.
func NewMemoryStorage() Storage {
	return &memoryStorage{
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}
}

// RestoreFileWithMetrics the function opens a file and saves metrics to it.
func (m *memoryStorage) RestoreFileWithMetrics(filename string) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("Сan not open file: %s\n", err)
	}
	defer file.Close()

	var allM []models.Metrics

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&allM)
	if err != nil {
		fmt.Printf("Сan not restore data: %s\n", err)
	}

	for _, metric := range allM {
		if metric.MType == metrics.Gauge {
			m.AddMetric(context.TODO(), metric.MType, *metric.Value, metric.ID)
		} else if metric.MType == metrics.Counter {
			m.AddMetric(context.TODO(), metric.MType, *metric.Delta, metric.ID)
		}
	}
}

// GetAllMetrics the function tries to get all the metrics contained in the memory.
func (m *memoryStorage) GetAllMetrics(_ context.Context, size int64) ([]models.Metrics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	allM := make([]models.Metrics, 0, size)
	for name, value := range m.gaugeMetrics {
		v := value
		m := models.Metrics{ID: name, MType: metrics.Gauge, Value: &v}
		allM = append(allM, m)
	}
	for name, value := range m.counterMetrics {
		v := value
		m := models.Metrics{ID: name, MType: metrics.Counter, Delta: &v}
		allM = append(allM, m)
	}
	return allM, nil
}

// GetCountMetrics the function tries to get count of all metrics contained in the memory.
func (m *memoryStorage) GetCountMetrics(_ context.Context) (int64, error) {
	s := len(m.gaugeMetrics) + len(m.counterMetrics)
	return int64(s), nil
}

// GetValueMetric the function tries to get value of concrete metric by type and name contained in the memory.
func (m *memoryStorage) GetValueMetric(_ context.Context, typeM string, name string) (interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var value interface{}
	var ok bool
	var err error
	if typeM == metrics.Gauge {
		value, ok = m.gaugeMetrics[name]
	}
	if typeM == metrics.Counter {
		value, ok = m.counterMetrics[name]
	}
	if !ok {
		err = fmt.Errorf("could not find metric with name %s", name)
		return nil, err
	}
	return value, nil
}

// GetAllMetricsJSON the function tries to get all the metrics in json format contained in the memory.
func (m *memoryStorage) GetAllMetricsJSON() []models.Metrics {
	m.mu.Lock()
	defer m.mu.Unlock()

	allM := make([]models.Metrics, 0, len(m.gaugeMetrics)+len(m.counterMetrics))
	for name, value := range m.gaugeMetrics {
		v := value
		m := models.Metrics{ID: name, MType: metrics.Gauge, Value: &v}
		allM = append(allM, m)
	}
	for name, value := range m.counterMetrics {
		v := value
		m := models.Metrics{ID: name, MType: metrics.Counter, Delta: &v}
		allM = append(allM, m)
	}
	return allM
}

// AddMetric the function tries to add value of concrete metric to memory.
func (m *memoryStorage) AddMetric(_ context.Context, typeM string, value interface{}, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if typeM == metrics.Counter {
		var newValue int64

		if existingValue, ok := m.counterMetrics[name]; ok {
			newValue = existingValue + value.(int64)
		} else {
			newValue = value.(int64)
		}
		m.counterMetrics[name] = newValue
	}

	if typeM == metrics.Gauge {
		m.gaugeMetrics[name] = value.(float64)
	}

	return nil
}

// WriteAndSaveMetricsToFile the function writes metrics to a file and saves them there.
func (m *memoryStorage) WriteAndSaveMetricsToFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	metrics := m.GetAllMetricsJSON()

	data, err := json.MarshalIndent(metrics, "", "   ")
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (m *memoryStorage) Ping() error {
	return nil
}
