package main

import (
	"github.com/kholodmv/go-service.git/cmd/handlers"
	"github.com/kholodmv/go-service.git/cmd/storage"
	"net/http"
)

func main() {
	metricRepository := storage.NewMetricRepository()
	metricHandler := handlers.NewMetricHandler(metricRepository)

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", metricHandler.UpdateMetric)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
	/*store := NewMetricsStore()

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", store.PostHandler)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}*/
}
