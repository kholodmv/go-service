// main the agent module is designed to collect metrics from various sources
// and transferring them to storage.
package main

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/kholodmv/go-service/cmd/metrics"
	"github.com/kholodmv/go-service/internal/configs"
	"github.com/kholodmv/go-service/internal/middleware/logger"
	"github.com/kholodmv/go-service/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
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

	conn, err := dialGRPC(conf)
	if err != nil {
		log.Fatal("unable to open client connection", zap.Error(err))
	}
	defer conn.Close()
	metricClient := proto.NewMetricsClient(conn)

	stop := make(chan os.Signal, 2)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx, stopMonitor := context.WithCancel(context.Background())
	connectionsClosed := make(chan struct{})

	m := metrics.Metrics{}
	m.ReportAgent(ctx, connectionsClosed, conf, publicKey, metricClient)

	log.Info("signal received, shutting down", zap.Any("signal", <-stop))
	stopMonitor()
	<-connectionsClosed

	log.Info("Shutting down agent")
}

func loadTLSCredentials(cafile string) (credentials.TransportCredentials, error) {
	pemServerCA, err := os.ReadFile(cafile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}

func dialGRPC(cfg configs.ConfigAgent) (*grpc.ClientConn, error) {
	creds, err := loadTLSCredentials(cfg.CACertFile)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(cfg.GRPCAddress, grpc.WithTransportCredentials(creds), grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
