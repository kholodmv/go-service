// main the agent module is designed to collect metrics from various sources
// and transferring them to storage.
package main

import (
	"github.com/kholodmv/go-service/cmd/metrics"
	"github.com/kholodmv/go-service/internal/configs"
)

func main() {
	conf := configs.InitConfigAgent()

	m := metrics.Metrics{}
	m.ReportAgent(conf)
}
