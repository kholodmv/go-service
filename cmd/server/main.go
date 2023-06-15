package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	get_all "github.com/kholodmv/go-service/cmd/handlers/getall"
	get_value "github.com/kholodmv/go-service/cmd/handlers/getvalue"
	"github.com/kholodmv/go-service/cmd/handlers/update"
	"github.com/kholodmv/go-service/cmd/storage"
	"net/http"
	"os"
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
	parseFlags()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	fmt.Println("Running server on", flagRunAddr)
	return http.ListenAndServe(flagRunAddr, MetricRouter())
}
