package configs

import (
	"flag"
	"fmt"
	"os"
)

type ServerConfig struct {
	RunAddress    string
	LogLevel      string
	StoreInterval int
	FileName      string
	Restore       bool
}

func UseServerStartParams() ServerConfig {
	var c ServerConfig

	flag.StringVar(&c.RunAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&c.LogLevel, "l", "info", "log level")
	flag.IntVar(&c.StoreInterval, "i", 20, "time interval in sec")
	flag.StringVar(&c.FileName, "f", "/tmp/metrics-db.json", "full file path")
	flag.BoolVar(&c.Restore, "r", true, "is load previously saved values")

	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		c.RunAddress = envRunAddr
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		c.LogLevel = envLogLevel
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		fmt.Sscanf(envStoreInterval, "%d", c.StoreInterval)
	}
	if envFileName := os.Getenv("FILE_STORAGE_PATH"); envFileName != "" {
		c.FileName = envFileName
	}
	if envFlagRestore := os.Getenv("RESTORE"); envFlagRestore != "" {
		fmt.Sscan(envFlagRestore, c.Restore)
	}

	return c
}
