package get_all

import (
	"fmt"
	"github.com/kholodmv/go-service/cmd/storage"
	"net/http"
)

type MetricHandler struct {
	metricStorage storage.MetricRepository
}

func NewMetricHandler(metricStorage storage.MetricRepository) *MetricHandler {
	return &MetricHandler{
		metricStorage: metricStorage,
	}
}

func (m *MetricHandler) GetAllMetric(res http.ResponseWriter, req *http.Request) {
	checkHTTPMethod(res, req)

	metrics := m.metricStorage.GetAllMetrics()

	var str string
	for _, metric := range metrics {
		str += fmt.Sprintf("%q : %v\n", metric.Name, metric.Value)
	}

	fmt.Fprint(res, str)
}

func checkHTTPMethod(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET methods", http.StatusMethodNotAllowed)
		return
	}
	res.Header().Set("Content-Type", "text/plain")
}
