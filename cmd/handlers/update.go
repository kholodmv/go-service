package handlers

import (
	"github.com/kholodmv/go-service/cmd/metrics"
	"github.com/kholodmv/go-service/cmd/storage"
	"net/http"
	"strconv"
	"strings"
)

type UpdateHandler struct {
	repository storage.MetricRepository
}

func NewUpdateHandler(repository storage.MetricRepository) *UpdateHandler {
	return &UpdateHandler{
		repository: repository,
	}
}

type PathParam struct {
	typeP string
	value string
	name  string
}

func (mh *UpdateHandler) UpdateMetric(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain")

	params := parsParams(res, req)

	checkType(res, params.typeP)
	checkName(res, params.name)

	checkAndSaveMetric(params.typeP, params.name, params.value, res, mh)

	res.WriteHeader(http.StatusOK)
}

func parsParams(res http.ResponseWriter, req *http.Request) PathParam {
	parts := strings.Split(req.URL.Path, "/")
	if len(parts) != 5 {
		http.Error(res, "Invalid request", http.StatusNotFound)
		return PathParam{}
	}
	p := PathParam{typeP: parts[2], name: parts[3], value: parts[4]}
	return p
}

func checkType(res http.ResponseWriter, metricType string) {
	if metricType != metrics.Gauge &&
		metricType != metrics.Counter {
		http.Error(res, "Incorrect type of metric "+metricType, http.StatusBadRequest)
		return
	}
}

func checkAndSaveMetric(metricType string, name string, value string, res http.ResponseWriter, mh *UpdateHandler) {
	switch metricType {
	case metrics.Gauge:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(res, "Invalid metric value", http.StatusBadRequest)
			return
		}
		mh.repository.AddGauge(v, name)

	case metrics.Counter:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(res, "Invalid metric value", http.StatusBadRequest)
		}
		mh.repository.AddCounter(v, name)
	}
}

func checkName(res http.ResponseWriter, metricName string) {
	if metricName == "" {
		http.Error(res, "Invalid metric name", http.StatusNotFound)
		return
	}
}
