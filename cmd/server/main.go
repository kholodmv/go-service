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

	handler := handlers.NewHandler(memoryStorage)
	handler.RegisterRoutes(router)

	return router
}

func main() {
	flags := configs.UseServerStartParams()

	fmt.Println("Running server on", flags)
	
	err := http.ListenAndServe(flags, MetricRouter())
	if err != nil {
		panic(err)
	}
}
