package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/kholodmv/go-service/cmd/handlers"
	"github.com/kholodmv/go-service/cmd/storage"
	"github.com/kholodmv/go-service/internal/configs"
	"github.com/kholodmv/go-service/internal/logger"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func MetricRouter(mem storage.MetricRepository) chi.Router {
	router := chi.NewRouter()

	handler := handlers.NewHandler(router, mem, configs.FlagFileName, configs.FlagRestore)
	handler.RegisterRoutes(router)

	return router
}

func main() {
	configs.UseServerStartParams()

	memoryStorage := storage.NewMemoryStorage()

	server := http.Server{
		Addr:    configs.FlagRunAddr,
		Handler: MetricRouter(memoryStorage),
	}

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(configs.FlagStoreInterval))
			memoryStorage.WriteAndSaveMetricsToFile(configs.FlagFileName)
		}
	}()

	connectionsClosed := make(chan struct{})
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		log.Println("Shutting down server")

		if err := memoryStorage.WriteAndSaveMetricsToFile(configs.FlagFileName); err != nil {
			log.Printf("Error during saving data to file: %v", err)
		}
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(connectionsClosed)
	}()

	logger.Initialize(configs.FlagLogLevel)
	logger.Log.Info("Running server", zap.String("address", configs.FlagRunAddr))

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
