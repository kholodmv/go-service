package storage

import (
	"encoding/json"
	"fmt"
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
	GetValueGauge(name string) (float64, error)
	GetValueCounterMetric(name string) (int64, bool)
	GetValueCounter(name string) (int64, error)
	GetAllMetrics() []Metric
	WriteAndSaveMetricsToFile(filename string) error
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

func (m *memoryStorage) GetValueGauge(name string) (float64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	value, ok := m.gaugeMetrics[name]
	if !ok {
		return value, fmt.Errorf("gauge metric with name '%s' not found", name)
	}
	return value, nil
}

func (m *memoryStorage) GetValueCounter(name string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	value, ok := m.counterMetrics[name]
	if !ok {
		return value, fmt.Errorf("counter metric with name '%s' not found", name)
	}
	return value, nil
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

	dataGauge, err := json.MarshalIndent(m.gaugeMetrics, "", "   ")
	if err != nil {
		return err
	}

	_, err = file.Write(dataGauge)
	if err != nil {
		return err
	}

	dataCounter, err := json.MarshalIndent(m.counterMetrics, "", "   ")
	if err != nil {
		return err
	}

	_, err = file.Write(dataCounter)
	if err != nil {
		return err
	}
	return nil
}
