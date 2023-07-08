package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/kholodmv/go-service/cmd/handlers"
	"github.com/kholodmv/go-service/cmd/storage"
	"github.com/kholodmv/go-service/internal/configs"
	"github.com/kholodmv/go-service/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

func MetricRouter() chi.Router {
	router := chi.NewRouter()

	memoryStorage := storage.NewMemoryStorage()

	handler := handlers.NewHandler(router, memoryStorage)
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

	logger.Log.Info("Running server", zap.String("address", configs.FlagRunAddr))
	return http.ListenAndServe(configs.FlagRunAddr, MetricRouter())
}
