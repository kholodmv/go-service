package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/kholodmv/go-service/cmd/handlers"
	"github.com/kholodmv/go-service/internal/configs"
	"github.com/kholodmv/go-service/internal/logger"
	"github.com/kholodmv/go-service/internal/store"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := configs.UseServerStartParams()
	log := logger.Initialize()

	var db store.Storage

	if cfg.DB != "" {
		log.Infow("DB main")
		db = store.NewStorage(cfg.DB)
	}

	router := chi.NewRouter()

	if cfg.Restore {
		db.RestoreFileWithMetrics(cfg.FileName)
	}

	handler := handlers.NewHandler(router, db, *log)
	handler.RegisterRoutes(router)

	server := http.Server{
		Addr:    cfg.RunAddress,
		Handler: router,
	}

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(cfg.StoreInterval))
			db.WriteAndSaveMetricsToFile(cfg.FileName)
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
		log.Info("Shutting down server")

		if err := db.WriteAndSaveMetricsToFile(cfg.FileName); err != nil {
			log.Errorf("Error during saving data to file: %v", err)
		}
		if err := server.Shutdown(context.Background()); err != nil {
			log.Errorf("HTTP Server Shutdown Error: %v", err)
		}
		close(connectionsClosed)
	}()

	log.Infow("Running server", zap.String("address", cfg.RunAddress))
	if err := server.ListenAndServe(); err != nil {
		log.Fatalw(err.Error(), "event", "start server")
	}
}
