package get_value

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kholodmv/go-service/cmd/storage"
	"io"
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

func (mh *MetricHandler) GetValueMetric(res http.ResponseWriter, req *http.Request) {
	checkHTTPMethod(res, req)

	typeMetric := chi.URLParam(req, "type")
	name := chi.URLParam(req, "name")

	var value interface{}
	var ok bool

	if typeMetric == "gauge" {
		value, ok = mh.metricStorage.GetValueGaugeMetric(name)
	}
	if typeMetric == "counter" {
		value, ok = mh.metricStorage.GetValueCounterMetric(name)
	}
	
	fmt.Println(value)
	if !ok {
		http.NotFound(res, req)
		return
	}
	strValue := fmt.Sprintf("%v", value)

	io.WriteString(res, strValue)
	res.WriteHeader(http.StatusOK)
}

func checkHTTPMethod(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET methods", http.StatusMethodNotAllowed)
		return
	}
	res.Header().Set("Content-Type", "text/plain")
}
