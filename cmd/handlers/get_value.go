package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kholodmv/go-service/cmd/metrics"
	"github.com/kholodmv/go-service/cmd/storage"
	"io"
	"net/http"
)

type GetValueHandler struct {
	repository storage.MetricRepository
}

func NewGetValueHandler(repository storage.MetricRepository) *GetValueHandler {
	return &GetValueHandler{
		repository: repository,
	}
}

func (mh *GetValueHandler) GetValueMetric(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain")

	typeMetric := chi.URLParam(req, "type")
	name := chi.URLParam(req, "name")

	var value interface{}
	var ok bool

	if typeMetric == metrics.Gauge {
		value, ok = mh.repository.GetValueGaugeMetric(name)
	}
	if typeMetric == metrics.Counter {
		value, ok = mh.repository.GetValueCounterMetric(name)
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
