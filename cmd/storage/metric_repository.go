package storage

import (
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
}

type memoryStorage struct {
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
	sync.Mutex
}

func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}
}

func (m *memoryStorage) GetValueGaugeMetric(name string) (float64, bool) {
	value, ok := m.gaugeMetrics[name]
	return value, ok
}

func (m *memoryStorage) GetValueCounterMetric(name string) (int64, bool) {
	value, ok := m.counterMetrics[name]
	return value, ok
}

func (m *memoryStorage) GetAllMetrics() []Metric {
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
	m.Lock()

	var newValue int64

	if existingValue, ok := m.counterMetrics[name]; ok {
		newValue = existingValue + value
	} else {
		newValue = value
	}
	m.counterMetrics[name] = newValue

	m.Unlock()
}
func (m *memoryStorage) AddGauge(value float64, name string) {
	m.Lock()

	m.gaugeMetrics[name] = value

	m.Unlock()
}
