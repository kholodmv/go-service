package main

import (
	"github.com/kholodmv/go-service/cmd/handlers"
	"github.com/kholodmv/go-service/cmd/storage"
	"net/http"
)

func main() {
	memoryStorage := storage.NewMetricMemoryStorage()
	handler := handlers.NewMetricHandler(memoryStorage)

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handler.UpdateMetric)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
