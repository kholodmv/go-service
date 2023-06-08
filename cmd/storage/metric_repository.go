package storage

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

type metricMemoryRepository struct {
	metrics map[string]Metric
}

func NewMetricRepository() *metricMemoryRepository {
	return &metricMemoryRepository{
		metrics: make(map[string]Metric),
	}
}

func (m *metricMemoryRepository) TypeCounter(value int64) {
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
}
func (m *metricMemoryRepository) TypeGauge(value float64) {
	var mm = Metric{
		Type:  Gauge,
		Value: value,
	}
	m.metrics[Gauge] = mm
}
