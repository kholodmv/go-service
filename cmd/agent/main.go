// main the agent module is designed to collect metrics from various sources
// and transferring them to storage.
package main

import (
	"context"
	"crypto/rsa"
	"github.com/kholodmv/go-service/cmd/metrics"
	"github.com/kholodmv/go-service/internal/configs"
	"github.com/kholodmv/go-service/internal/middleware/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

var (
	buildVersion, buildDate, buildCommit string = "N/A", "N/A", "N/A"
)

func main() {
	log := logger.Initialize()
	defer log.Sync()
	log.Infof("\nBuild version: %v", buildVersion)
	log.Infof("\nBuild date: %v", buildDate)
	log.Infof("\nBuild commit: %v", buildCommit)

	conf := configs.InitConfigAgent()

	var publicKey *rsa.PublicKey
	var err error
	if conf.CryptoPublicKey != "" {
		publicKey, err = conf.GetPublicKey()
		if err != nil {
			log.Error("failed to get public encryption key", zap.Error(err))
			return
		}
	}

	if conf.ConfigFile != "" {
		if err = conf.ParseFile(conf.ConfigFile); err != nil {
			log.Error("failed to parse file", zap.Error(err))
		}
	}

	stop := make(chan os.Signal, 2)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx, stopMonitor := context.WithCancel(context.Background())
	connectionsClosed := make(chan struct{})

	m := metrics.Metrics{}
	m.ReportAgent(ctx, connectionsClosed, conf, publicKey)

	log.Info("signal received, shutting down", zap.Any("signal", <-stop))
	stopMonitor()
	<-connectionsClosed

	log.Info("Shutting down agent")
}
