package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kholodmv/go-service/cmd/metrics"
	"github.com/kholodmv/go-service/internal/models"
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

func (mh *Handler) UpdateJsonMetric(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	var m models.Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &m); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	switch m.MType {
	case metrics.Counter:
		if m.Delta == nil {
			http.Error(res, "Metric value type counter should not be empty", http.StatusBadRequest)
			return
		}
		mh.repository.AddCounter(*m.Delta, m.ID)
		res.WriteHeader(http.StatusOK)
	case metrics.Gauge:
		if m.Value == nil {
			http.Error(res, "Metric value type gauge should not be empty", http.StatusBadRequest)
			return
		}
		mh.repository.AddGauge(*m.Value, m.ID)
		res.WriteHeader(http.StatusOK)
	default:
		http.Error(res, "Incorrect metric type", http.StatusBadRequest)
	}
}

func (mh *Handler) GetJsonMetric(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	var m models.Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &m); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	switch m.MType {
	case metrics.Counter:
		counter, ok := mh.repository.GetValueCounterMetric(m.ID)
		if !ok {
			http.NotFound(res, req)
			return
		}
		m.Delta = &counter
		m.MType = metrics.Counter
	case metrics.Gauge:
		gauge, ok := mh.repository.GetValueGaugeMetric(m.ID)
		if !ok {
			http.NotFound(res, req)
			return
		}
		m.Value = &gauge
		m.MType = metrics.Gauge
	}

	resp, err := json.Marshal(m)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(resp)
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
	res.WriteHeader(http.StatusOK)
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
