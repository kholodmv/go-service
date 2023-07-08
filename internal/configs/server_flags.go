package configs

import (
	"flag"
	"fmt"
	"os"
)

var (
	FlagRunAddr       string
	FlagLogLevel      string
	FlagStoreInterval int
	FlagFileName      string
	FlagRestore       bool
)

func UseServerStartParams() {
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&FlagLogLevel, "l", "info", "log level")
	flag.IntVar(&FlagStoreInterval, "i", 10, "time interval in sec")
	flag.StringVar(&FlagFileName, "f", "/tmp/metrics-db.json", "full file path")
	flag.BoolVar(&FlagRestore, "r", true, "is load previously saved values")

	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		FlagLogLevel = envLogLevel
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		fmt.Sscanf(envStoreInterval, "%d", FlagStoreInterval)

	}
	if envFileName := os.Getenv("FILE_STORAGE_PATH"); envFileName != "" {
		FlagFileName = envFileName
	}
	if envFlagRestore := os.Getenv("RESTORE"); envFlagRestore != "" {
		fmt.Sscan(envFlagRestore, FlagRestore)
	}
}
