// main the agent module is designed to collect metrics from various sources
// and transferring them to storage.
package main

import (
	"github.com/kholodmv/go-service/cmd/metrics"
	"github.com/kholodmv/go-service/internal/configs"
	"github.com/kholodmv/go-service/internal/logger"
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

	m := metrics.Metrics{}
	m.ReportAgent(conf)
}
