package handlers

import (
	"github.com/kholodmv/go-service.git/cmd/storage"
	"net/http"
	"strconv"
	"strings"
)

const (
	Gauge   string = "gauge"
	Counter string = "counter"
)

type MetricHandler struct {
	metricRepository storage.MetricRepository
}

func NewMetricHandler(metricRepository storage.MetricRepository) *MetricHandler {
	return &MetricHandler{
		metricRepository: metricRepository,
	}
}

func (uh *MetricHandler) UpdateMetric(res http.ResponseWriter, req *http.Request) {
	checkHTTPMethod(res, req)

	parts := strings.Split(req.URL.Path, "/")
	if len(parts) != 5 {
		http.Error(res, "Invalid request", http.StatusNotFound)
		return
	}

	checkType(res, parts, uh)

	metricName := parts[3]
	checkName(res, metricName)

	res.WriteHeader(http.StatusOK)
}

func checkType(res http.ResponseWriter, parts []string, uh *MetricHandler) {
	metricType := parts[2]

	switch metricType {
	case Gauge:
		value, err := strconv.ParseFloat(parts[4], 64)
		if err != nil {
			http.Error(res, "Invalid metric value", http.StatusBadRequest)
			return
		}
		uh.metricRepository.TypeGauge(value)

	case Counter:
		value, err := strconv.ParseInt(parts[4], 10, 64)
		if err != nil {
			http.Error(res, "Invalid metric value", http.StatusBadRequest)
		}
		uh.metricRepository.TypeCounter(value)

	default:
		http.Error(res, "Incorrect type of metric "+metricType, http.StatusBadRequest)
		return
	}
}

func checkName(res http.ResponseWriter, metricName string) {
	if metricName == "" {
		http.Error(res, "Invalid metric name", http.StatusNotFound)
		return
	}
}

func checkHTTPMethod(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST methods", http.StatusMethodNotAllowed)
		return
	}
	res.Header().Set("Content-Type", "application/json")
}
