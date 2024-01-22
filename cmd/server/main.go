// main the server module is entry point to the program.
package main

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"github.com/kholodmv/go-service/internal/core"
	"github.com/kholodmv/go-service/internal/interceptors"
	"github.com/kholodmv/go-service/internal/middleware/logger"
	"github.com/kholodmv/go-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
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

	var privateKey *rsa.PrivateKey
	var err error
	if cfg.CryptoPrivateKey != "" {
		privateKey, err = cfg.GetPrivateKey()
		if err != nil {
			log.Error("failed to get private encryption key", zap.Error(err))
			return
		}
	} else {
		privateKey = nil
	}

	if cfg.ConfigFile != "" {
		if err = cfg.ParseFile(cfg.ConfigFile); err != nil {
			log.Error("failed to parse file", zap.Error(err))
		}
	}

	handler := handlers.NewHandler(router, db, *log, cfg.Key, cfg.TrustedSubnet, privateKey)
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
	grpc := startGRPC(cfg, log, db)

	connectionsClosed := make(chan struct{})
	go func() {
		stop := make(chan os.Signal, 2)
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
	grpc.GracefulStop()
}

// connectToDB is function which connected to postgres db.
func connectToDB(path string) *sql.DB {
	con, err := sql.Open("postgres", path)
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	return con
}

func startGRPC(cfg configs.ServerConfig, log *zap.SugaredLogger, repo store.Storage) *grpc.Server {
	listen, err := net.Listen("tcp", cfg.GAddress)
	if err != nil {
		log.Fatal("unable to listen tcp", zap.Error(err))
	}

	opts := make([]grpc.ServerOption, 0)
	if cfg.GRPCConfig.TLSCertFile != "" && cfg.GRPCConfig.TLSKeyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(cfg.TLSCertFile, cfg.GRPCConfig.TLSKeyFile)
		if err != nil {
			log.Fatal("failed to create credentials: %v", zap.Error(err))
		}
		opts = append(opts, grpc.Creds(creds))
	}

	opts = append(opts, interceptors.RegisterUnaryInterceptorChain(cfg))

	s := grpc.NewServer(opts...)
	proto.RegisterMetricsServer(s, core.NewMetricsServer(repo, log.Desugar()))

	go func() {
		if err := s.Serve(listen); err != nil {
			log.Fatal("unable to start grpc server", zap.Error(err))
		}
	}()

	return s
}
