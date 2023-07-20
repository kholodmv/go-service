package configs

import (
	"flag"
	"fmt"
	"github.com/kholodmv/go-service/internal/logger"
	"go.uber.org/zap"
	"os"
)

type ServerConfig struct {
	DB            string
	RunAddress    string
	LogLevel      string
	StoreInterval int
	FileName      string
	Restore       bool
}

func UseServerStartParams() ServerConfig {
	var c ServerConfig
	log := logger.Initialize()

	flag.StringVar(&c.DB, "d", fmt.Sprintf("host=%s port=%d dbname=%s sslmode=disable",
		"localhost", 5000, "postgres"),
		"connection string to postgres db")
	flag.StringVar(&c.RunAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&c.LogLevel, "l", "info", "log level")
	flag.IntVar(&c.StoreInterval, "i", 300, "time interval in sec")
	flag.StringVar(&c.FileName, "f", "/tmp/metrics-db.json", "full file path")
	flag.BoolVar(&c.Restore, "r", true, "is load previously saved values")

	flag.Parse()

	log.Infow("1 Conf log", zap.String("db", c.DB))

	if envRunDB := os.Getenv("DATABASE_DSN"); envRunDB != "" {
		c.DB = envRunDB
	}
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

	log.Infow("2 Conf log", zap.String("db", c.DB))

	return c
}
