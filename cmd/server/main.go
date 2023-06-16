package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	get_all "github.com/kholodmv/go-service/cmd/handlers/getall"
	get_value "github.com/kholodmv/go-service/cmd/handlers/getvalue"
	"github.com/kholodmv/go-service/cmd/handlers/update"
	"github.com/kholodmv/go-service/cmd/storage"
	"net/http"
)

func MetricRouter() chi.Router {
	router := chi.NewRouter()

	memoryStorage := storage.NewMemoryStorage()

	updHandler := update.NewHandler(memoryStorage)
	getValueHandler := get_value.NewHandler(memoryStorage)
	getAllHandler := get_all.NewHandler(memoryStorage)

	router.Post("/update/{type}/{name}/{value}", updHandler.UpdateMetric)
	router.Get("/value/{type}/{name}", getValueHandler.GetValueMetric)
	router.Get("/", getAllHandler.GetAllMetric)

	return router
}

func main() {
	flags := useStartParams()

	if err := run(flags); err != nil {
		panic(err)
	}
}

func run(flags string) error {
	fmt.Println("Running server on", flags)
	return http.ListenAndServe(flags, MetricRouter())
}
