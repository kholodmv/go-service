package storage

import "sync"

const (
	Gauge   string = "gauge"
	Counter string = "counter"
)

type Metric struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type MetricRepository interface {
	TypeCounter(value int64)
	TypeGauge(value float64)
}

type metricMemoryStorage struct {
	metrics map[string]Metric
	sync.Mutex
}

func NewMetricMemoryStorage() *metricMemoryStorage {
	return &metricMemoryStorage{
		metrics: make(map[string]Metric),
	}
}

func (m *metricMemoryStorage) TypeCounter(value int64) {
	m.Lock()

	var mm = Metric{}
	if existingValue, ok := m.metrics[Counter].Value.(int64); ok {
		mm = Metric{
			Type:  Gauge,
			Value: existingValue + value,
		}
	} else {
		mm = Metric{
			Type:  Gauge,
			Value: value,
		}
	}
	m.metrics[Counter] = mm

	m.Unlock()
}
func (m *metricMemoryStorage) TypeGauge(value float64) {
	m.Lock()

	var mm = Metric{
		Type:  Gauge,
		Value: value,
	}
	m.metrics[Gauge] = mm

	m.Unlock()
}
