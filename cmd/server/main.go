package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/kholodmv/go-service/cmd/handlers"
	"github.com/kholodmv/go-service/cmd/storage"
	"github.com/kholodmv/go-service/internal/configs"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := configs.UseServerStartParams()
	memoryStorage := storage.NewMemoryStorage()
	router := chi.NewRouter()

	handler := handlers.NewHandler(router, memoryStorage, cfg.FileName, cfg.Restore)
	handler.RegisterRoutes(router)

	server := http.Server{
		Addr:    cfg.RunAddress,
		Handler: router,
	}

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(cfg.StoreInterval))
			memoryStorage.WriteAndSaveMetricsToFile(cfg.FileName)
		}
	}()

	connectionsClosed := make(chan struct{})
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)
		<-stop
		log.Println("Shutting down server")

		if err := memoryStorage.WriteAndSaveMetricsToFile(cfg.FileName); err != nil {
			log.Printf("Error during saving data to file: %v", err)
		}
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(connectionsClosed)
	}()

	log.Println("Running server", zap.String("address", cfg.RunAddress))
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
