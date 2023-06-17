package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kholodmv/go-service/cmd/handlers"
	"github.com/kholodmv/go-service/cmd/storage"
	"github.com/kholodmv/go-service/internal/configs"
	"net/http"
)

func MetricRouter() chi.Router {
	router := chi.NewRouter()

	memoryStorage := storage.NewMemoryStorage()

	updHandler := handlers.NewUpdateHandler(memoryStorage)
	getValueHandler := handlers.NewGetValueHandler(memoryStorage)
	getAllHandler := handlers.NewGetAllHandler(memoryStorage)

	router.Post("/update/{type}/{name}/{value}", updHandler.UpdateMetric)
	router.Get("/value/{type}/{name}", getValueHandler.GetValueMetric)
	router.Get("/", getAllHandler.GetAllMetric)

	return router
}

func main() {
	flags := configs.UseServerStartParams()

	if err := run(flags); err != nil {
		panic(err)
	}
}

func run(flags string) error {
	fmt.Println("Running server on", flags)
	return http.ListenAndServe(flags, MetricRouter())
}
