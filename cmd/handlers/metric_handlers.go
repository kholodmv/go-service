package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kholodmv/go-service/cmd/metrics"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type PathParam struct {
	typeP string
	value string
	name  string
}

func (mh *Handler) GetValueMetric(res http.ResponseWriter, req *http.Request) {
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

	if !ok {
		http.NotFound(res, req)
		return
	}
	strValue := fmt.Sprintf("%v", value)

	io.WriteString(res, strValue)
	res.WriteHeader(http.StatusOK)
}

func (mh *Handler) GetAllMetric(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain")

	metrics := mh.repository.GetAllMetrics()

	var str string
	for _, metric := range metrics {
		str += fmt.Sprintf("%q : %v\n", metric.Name, metric.Value)
	}

	fmt.Fprint(res, str)
}

func (mh *Handler) UpdateMetric(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain")

	params, err := isValidParams(req)
	if !err {
		http.Error(res, "Invalid request", http.StatusNotFound)
		return
	}

	err = checkType(params.typeP)
	if !err {
		http.Error(res, "Incorrect type of metric "+params.typeP, http.StatusBadRequest)
		return
	}

	err = checkName(params.name)
	if !err {
		http.Error(res, "Invalid metric name", http.StatusNotFound)
		return
	}

	err = checkAndSaveMetric(params.typeP, params.name, params.value, mh)
	if !err {
		http.Error(res, "Invalid metric value", http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func isValidParams(req *http.Request) (PathParam, bool) {
	parts := strings.Split(req.URL.Path, "/")
	if len(parts) != 5 {
		return PathParam{}, false
	}
	p := PathParam{typeP: parts[2], name: parts[3], value: parts[4]}
	return p, true
}

func checkType(metricType string) bool {
	if metricType != metrics.Gauge &&
		metricType != metrics.Counter {
		return false
	}
	return true
}

func checkAndSaveMetric(metricType string, name string, value string, mh *Handler) bool {
	switch metricType {
	case metrics.Gauge:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return false
		}
		mh.repository.AddGauge(v, name)

	case metrics.Counter:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return false
		}
		mh.repository.AddCounter(v, name)
	}
	return true
}

func checkName(metricName string) bool {
	return metricName != ""
}
