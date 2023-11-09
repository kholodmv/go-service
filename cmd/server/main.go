// main the server module is entry point to the program.
package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/kholodmv/go-service/cmd/handlers"
	"github.com/kholodmv/go-service/internal/configs"
	"github.com/kholodmv/go-service/internal/logger"
	"github.com/kholodmv/go-service/internal/store"
)

var (
	buildVersion, buildDate, buildCommit string = "N/A", "N/A", "N/A"
)

func main() {
	cfg := configs.UseServerStartParams()
	log := logger.Initialize()

	log.Infof("\nBuild version: %v", buildVersion)
	log.Infof("\nBuild date: %v", buildDate)
	log.Infof("\nBuild commit: %v", buildCommit)

	var db store.Storage

	if cfg.DB != "" {
		con := connectToDB(cfg.DB)
		s := store.NewStorage(con, log)
		store.CreateTable(s)
		db = s

	} else {
		db = store.NewMemoryStorage()
	}

	router := chi.NewRouter()

	if cfg.Restore {
		db.RestoreFileWithMetrics(cfg.FileName)
	}

	handler := handlers.NewHandler(router, db, *log, cfg.Key)
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

// connectToDB is function which connected to postgres db.
func connectToDB(path string) *sql.DB {
	con, err := sql.Open("postgres", path)
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	return con
}
