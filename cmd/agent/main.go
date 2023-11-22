// main the agent module is designed to collect metrics from various sources
// and transferring them to storage.
package main

import (
	"crypto/rsa"
	"github.com/kholodmv/go-service/cmd/metrics"
	"github.com/kholodmv/go-service/internal/configs"
	"github.com/kholodmv/go-service/internal/logger"
	"go.uber.org/zap"
)

var (
	buildVersion, buildDate, buildCommit string = "N/A", "N/A", "N/A"
)

func main() {
	log := logger.Initialize()
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

	m := metrics.Metrics{}
	m.ReportAgent(conf, publicKey)
}
