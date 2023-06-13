package main

import (
	"github.com/go-chi/chi/v5"
	get_all "github.com/kholodmv/go-service/cmd/handlers/get-all"
	get_value "github.com/kholodmv/go-service/cmd/handlers/get-value"
	"github.com/kholodmv/go-service/cmd/handlers/update"
	"github.com/kholodmv/go-service/cmd/storage"
	"net/http"
)

func MetricRouter() chi.Router {
	r := chi.NewRouter()

	memoryStorage := storage.NewMemoryStorage()

	updHandler := update.NewMetricHandler(memoryStorage)
	getValueHandler := get_value.NewMetricHandler(memoryStorage)
	getAllHandler := get_all.NewMetricHandler(memoryStorage)

	r.Post("/update/{type}/{name}/{value}", updHandler.UpdateMetric)
	r.Get("/value/{type}/{name}", getValueHandler.GetValueMetric)
	r.Get("/", getAllHandler.GetAllMetric)

	return r
}

func main() {
	http.ListenAndServe(":8080", MetricRouter())
}
