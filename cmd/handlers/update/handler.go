package update

import (
	"github.com/kholodmv/go-service/cmd/storage"
	"net/http"
	"strconv"
	"strings"
)

const (
	Gauge   string = "gauge"
	Counter string = "counter"
)

type MetricHandler struct {
	metricStorage storage.MetricRepository
}

func NewMetricHandler(metricStorage storage.MetricRepository) *MetricHandler {
	return &MetricHandler{
		metricStorage: metricStorage,
	}
}

func (mh *MetricHandler) UpdateMetric(res http.ResponseWriter, req *http.Request) {
	checkHTTPMethod(res, req)

	parts := strings.Split(req.URL.Path, "/")
	if len(parts) != 5 {
		http.Error(res, "Invalid request", http.StatusNotFound)
		return
	}
	metricName := parts[3]
	checkName(res, metricName)

	checkType(res, parts, mh)

	res.WriteHeader(http.StatusOK)
}

func checkType(res http.ResponseWriter, parts []string, mh *MetricHandler) {
	metricType := parts[2]
	metricName := parts[3]

	switch metricType {
	case Gauge:
		value, err := strconv.ParseFloat(parts[4], 64)
		if err != nil {
			http.Error(res, "Invalid metric value", http.StatusBadRequest)
			return
		}

		mh.metricStorage.AddGauge(value, metricName)

	case Counter:
		value, err := strconv.ParseInt(parts[4], 10, 64)
		if err != nil {
			http.Error(res, "Invalid metric value", http.StatusBadRequest)
		}
		mh.metricStorage.AddCounter(value, metricName)

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
	res.Header().Set("Content-Type", "text/plain")
}
