package storage

import (
	"encoding/json"
	"fmt"
	"github.com/kholodmv/go-service/internal/models"
	"os"
	"sync"
)

type Metric struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type MetricRepository interface {
	AddCounter(value int64, name string)
	AddGauge(value float64, name string)
	GetValueGaugeMetric(name string) (float64, bool)
	GetValueCounterMetric(name string) (int64, bool)
	GetAllMetrics() []Metric
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
			m.AddGauge(*metric.Value, metric.ID)
		} else if metric.MType == "counter" {
			m.AddCounter(*metric.Delta, metric.ID)
		}
	}
}

func (m *memoryStorage) GetValueGaugeMetric(name string) (float64, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	value, ok := m.gaugeMetrics[name]
	return value, ok
}

func (m *memoryStorage) GetValueCounterMetric(name string) (int64, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	value, ok := m.counterMetrics[name]
	return value, ok
}

func (m *memoryStorage) GetAllMetrics() []Metric {
	m.mu.Lock()
	defer m.mu.Unlock()

	metrics := make([]Metric, 0, len(m.gaugeMetrics)+len(m.counterMetrics))
	for name, value := range m.gaugeMetrics {
		metrics = append(metrics, Metric{Name: name, Value: value})
	}
	for name, value := range m.counterMetrics {
		metrics = append(metrics, Metric{Name: name, Value: value})
	}
	return metrics
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

func (m *memoryStorage) AddCounter(value int64, name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var newValue int64

	if existingValue, ok := m.counterMetrics[name]; ok {
		newValue = existingValue + value
	} else {
		newValue = value
	}
	m.counterMetrics[name] = newValue
}
func (m *memoryStorage) AddGauge(value float64, name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.gaugeMetrics[name] = value
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
