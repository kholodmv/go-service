package main

import (
	"net/http"
	"strconv"
	"strings"
)

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

func (store *MemStorage) PostHandler(res http.ResponseWriter, req *http.Request) {
	checkHTTPMethod(res, req)

	res.Header().Set("Content-Type", "application/json")

	parts := strings.Split(req.URL.Path, "/")
	if len(parts) != 5 {
		http.Error(res, "Invalid request", http.StatusNotFound)
		return
	}

	//store.mutex.Lock()
	//defer store.mutex.Unlock()

	checkType(res, parts, store)

	metricName := parts[3]
	checkName(res, metricName)
}

func checkType(res http.ResponseWriter, parts []string, store *MemStorage) {
	metricType := MetricType(parts[2])

	switch metricType {
	case Gauge:
		value, err := strconv.ParseFloat(parts[4], 64)
		if err != nil {
			http.Error(res, "Invalid metric value", http.StatusBadRequest)
			return
		}
		store.metrics["gauge"] = value

	case Counter:
		value, err := strconv.ParseInt(parts[4], 10, 64)
		if err != nil {
			http.Error(res, "Invalid metric value", http.StatusBadRequest)
		}
		if existingValue, ok := store.metrics["counter"].(int64); ok {
			store.metrics["counter"] = existingValue + value
		} else {
			store.metrics["counter"] = value
		}

	default:
		http.Error(res, string("Incorrect type of metric "+metricType), http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
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
}
