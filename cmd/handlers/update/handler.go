package update

import (
	"github.com/kholodmv/go-service/cmd/common"
	"github.com/kholodmv/go-service/cmd/storage"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	repository storage.MetricRepository
}

func NewHandler(repository storage.MetricRepository) *Handler {
	return &Handler{
		repository: repository,
	}
}

func (mh *Handler) UpdateMetric(res http.ResponseWriter, req *http.Request) {
	common.CheckPostHTTPMethod(res, req)
	res.Header().Set("Content-Type", "text/plain")

	params := parsParams(res, req)

	metricType := params[2]
	checkType(res, metricType)

	metricName := params[3]
	checkName(res, metricName)

	metricValue := params[4]
	checkMetricsValue(res, metricValue, metricType)

	saveMetrics(metricType, metricName, metricValue, mh)

	res.WriteHeader(http.StatusOK)
}

func parsParams(res http.ResponseWriter, req *http.Request) []string {
	parts := strings.Split(req.URL.Path, "/")
	if len(parts) != 5 {
		http.Error(res, "Invalid request", http.StatusNotFound)
		return nil
	}
	return parts
}

func checkType(res http.ResponseWriter, metricType string) {
	if metricType != common.Gauge && metricType != common.Counter {
		http.Error(res, "Incorrect type of metric "+metricType, http.StatusBadRequest)
		return
	}
}

func checkMetricsValue(res http.ResponseWriter, value string, typeMetric string) {
	if typeMetric == common.Gauge {
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(res, "Invalid metric value", http.StatusBadRequest)
			return
		}
	}

	if typeMetric == common.Counter {
		_, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(res, "Invalid metric value", http.StatusBadRequest)
		}
	}
}

func saveMetrics(metricType string, name string, value string, mh *Handler) {
	switch metricType {
	case common.Gauge:
		v, _ := strconv.ParseFloat(value, 64)
		mh.repository.AddGauge(v, name)

	case common.Counter:
		v, _ := strconv.ParseInt(value, 10, 64)
		mh.repository.AddCounter(v, name)
	}
}

func checkName(res http.ResponseWriter, metricName string) {
	if metricName == "" {
		http.Error(res, "Invalid metric name", http.StatusNotFound)
		return
	}
}
