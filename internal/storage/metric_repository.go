package storage

import (
	"encoding/json"
	"fmt"
	"github.com/kholodmv/go-service/cmd/metrics"
	"github.com/kholodmv/go-service/internal/models"
	"os"
	"sync"
)

type MetricRepository interface {
	AddMetric(typeM string, value interface{}, name string)
	GetValueMetric(typeM string, name string) (interface{}, bool)
	GetAllMetrics() []models.Metrics
	GetAllMetricsJSON() []models.Metrics
	WriteAndSaveMetricsToFile(filename string) error
	RestoreFileWithMetrics(filename string)
}

type memoryStorage struct {
	mu             sync.Mutex
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
}

func NewMemoryStorage() MetricRepository {
	return &memoryStorage{
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}
}

func (m *memoryStorage) RestoreFileWithMetrics(filename string) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("Сan not open file: %s\n", err)
	}
	defer file.Close()

	var metrics []models.Metrics

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&metrics)
	if err != nil {
		fmt.Printf("Сan not restore data: %s\n", err)
	}

	for _, metric := range metrics {
		if metric.MType == "gauge" {
			m.AddMetric(metric.MType, *metric.Value, metric.ID)
		} else if metric.MType == "counter" {
			m.AddMetric(metric.MType, *metric.Delta, metric.ID)
		}
	}
}

func (m *memoryStorage) GetValueMetric(typeM string, name string) (interface{}, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var value interface{}
	var ok bool
	if typeM == metrics.Gauge {
		value, ok = m.gaugeMetrics[name]
	}
	if typeM == metrics.Counter {
		value, ok = m.counterMetrics[name]
	}
	return value, ok
}

func (m *memoryStorage) GetAllMetrics() []models.Metrics {
	m.mu.Lock()
	defer m.mu.Unlock()

	allM := make([]models.Metrics, 0, len(m.gaugeMetrics)+len(m.counterMetrics))
	for name, value := range m.gaugeMetrics {
		v := value
		allM = append(allM, models.Metrics{ID: name, MType: metrics.Gauge, Value: &v})
	}
	for name, value := range m.counterMetrics {
		v := value
		allM = append(allM, models.Metrics{ID: name, MType: metrics.Counter, Delta: &v})
	}
	return allM
}

func (m *memoryStorage) GetAllMetricsJSON() []models.Metrics {
	m.mu.Lock()
	defer m.mu.Unlock()

	metrics := make([]models.Metrics, 0, len(m.gaugeMetrics)+len(m.counterMetrics))
	for name, value := range m.gaugeMetrics {
		v := value
		m := models.Metrics{ID: name, MType: "gauge", Value: &v}
		metrics = append(metrics, m)
	}
	for name, value := range m.counterMetrics {
		v := value
		m := models.Metrics{ID: name, MType: "counter", Delta: &v}
		metrics = append(metrics, m)
	}
	return metrics
}

func (m *memoryStorage) AddMetric(typeM string, value interface{}, name string) {
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
}

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
