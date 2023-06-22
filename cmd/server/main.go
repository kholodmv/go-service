package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kholodmv/go-service/cmd/handlers"
	"github.com/kholodmv/go-service/cmd/storage"
	"github.com/kholodmv/go-service/internal/configs"
	"github.com/kholodmv/go-service/internal/logger"
	"net/http"
)

func MetricRouter() chi.Router {
	router := chi.NewRouter()

	memoryStorage := storage.NewMemoryStorage()

	handler := handlers.NewHandler(memoryStorage)
	handler.RegisterRoutes(router)

	return router
}

func main() {
	configs.UseServerStartParams()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if err := logger.Initialize(configs.FlagLogLevel); err != nil {
		return err
	}

	fmt.Println("Running server on", configs.FlagRunAddr)
	return http.ListenAndServe(configs.FlagRunAddr, MetricRouter())
}
