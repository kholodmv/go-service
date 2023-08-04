package handlers

import (
	"bytes"
	"context"
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

func (mh *Handler) DBConnection(res http.ResponseWriter, _ *http.Request) {
	err := mh.db.Ping()
	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("OK"))
}

func (mh *Handler) UpdateJSONMetric(res http.ResponseWriter, req *http.Request) {
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
		mh.db.AddMetric(req.Context(), metrics.Counter, *m.Delta, m.ID)
		res.WriteHeader(http.StatusOK)

	case metrics.Gauge:
		if m.Value == nil {
			http.Error(res, "Metric value type gauge should not be empty", http.StatusBadRequest)
			return
		}
		mh.db.AddMetric(req.Context(), metrics.Gauge, *m.Value, m.ID)
		res.WriteHeader(http.StatusOK)

	default:
		http.Error(res, "Incorrect metric type", http.StatusBadRequest)
	}

	resp, err := json.Marshal(m)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Write(resp)
}

func (mh *Handler) UpdatesMetrics(res http.ResponseWriter, req *http.Request) {
	var m []models.Metrics
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

	for _, v := range m {
		switch v.MType {
		case metrics.Counter:
			if v.Delta == nil {
				http.Error(res, "Metric value type counter should not be empty", http.StatusBadRequest)
				return
			}
			mh.db.AddMetric(req.Context(), metrics.Counter, *v.Delta, v.ID)
			res.WriteHeader(http.StatusOK)

		case metrics.Gauge:
			if v.Value == nil {
				http.Error(res, "Metric value type gauge should not be empty", http.StatusBadRequest)
				return
			}
			mh.db.AddMetric(req.Context(), metrics.Gauge, *v.Value, v.ID)
			res.WriteHeader(http.StatusOK)

		default:
			http.Error(res, "Incorrect metric type", http.StatusBadRequest)
		}
	}

	res.WriteHeader(http.StatusOK)
}

func (mh *Handler) GetJSONMetric(res http.ResponseWriter, req *http.Request) {
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
		counter, err := mh.db.GetValueMetric(req.Context(), metrics.Counter, m.ID)
		if err != nil {
			http.NotFound(res, req)
			return
		}
		v := counter.(int64)
		m.Delta = &v
		m.MType = metrics.Counter

	case metrics.Gauge:
		gauge, err := mh.db.GetValueMetric(req.Context(), metrics.Gauge, m.ID)
		if err != nil {
			http.NotFound(res, req)
			return
		}
		v := gauge.(float64)
		m.Value = &v
		m.MType = metrics.Gauge
	}

	resp, err := json.Marshal(m)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}

func (mh *Handler) GetValueMetric(res http.ResponseWriter, req *http.Request) {
	if mh.key == "" {
		res.Header().Set("Content-Type", "text/plain")
	}

	typeMetric := chi.URLParam(req, "type")
	name := chi.URLParam(req, "name")

	var value interface{}

	if typeMetric == metrics.Gauge {
		value, _ = mh.db.GetValueMetric(req.Context(), metrics.Gauge, name)
		if value == nil {
			http.NotFound(res, req)
			return
		}
	} else if typeMetric == metrics.Counter {
		value, _ = mh.db.GetValueMetric(req.Context(), metrics.Counter, name)
		if value == nil {
			http.NotFound(res, req)
			return
		}
	} else {
		http.NotFound(res, req)
		return
	}

	strValue := fmt.Sprintf("%v", value)

	res.Header().Set("Content-Type", "application/json")
	io.WriteString(res, strValue)
	res.WriteHeader(http.StatusOK)
}

func (mh *Handler) GetAllMetric(res http.ResponseWriter, req *http.Request) {
	if mh.key != "" {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
	}

	size, _ := mh.db.GetCountMetrics(req.Context())
	allM, _ := mh.db.GetAllMetrics(req.Context(), size)

	var str string

	for _, metric := range allM {
		if metric.MType == metrics.Gauge {
			v := strconv.FormatFloat(*metric.Value, 'g', 5, 64)
			str += fmt.Sprintf("%q : %s\n", metric.ID, v)
		}
		if metric.MType == metrics.Counter {
			v := strconv.FormatInt(*metric.Delta, 10)
			str += fmt.Sprintf("%q : %s\n", metric.ID, v)
		}
	}

	fmt.Fprint(res, str)
	if mh.key != "" {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	res.WriteHeader(http.StatusOK)
}

func (mh *Handler) UpdateMetric(res http.ResponseWriter, req *http.Request) {
	if mh.key == "" {
		res.Header().Set("Content-Type", "text/plain")
	}

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
	err = checkAndSaveMetric(params.typeP, params.name, params.value, mh, req.Context())
	if !err {
		http.Error(res, "Invalid metric value", http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
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

func checkAndSaveMetric(metricType string, name string, value string, mh *Handler, ctx context.Context) bool {
	switch metricType {
	case metrics.Gauge:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return false
		}
		mh.db.AddMetric(ctx, metrics.Gauge, v, name)

	case metrics.Counter:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return false
		}
		mh.db.AddMetric(ctx, metrics.Counter, v, name)
	}
	return true
}

func checkName(metricName string) bool {
	return metricName != ""
}
