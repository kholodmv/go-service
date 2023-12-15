// Package configs - server_flags.go - server configuration parameters.
package configs

import (
	"flag"
	"fmt"
	"os"
)

// ServerConfig structure that contains variables for initial.
type ServerConfig struct {
	GRPCConfig
	DB               string
	RunAddress       string
	LogLevel         string
	StoreInterval    int
	FileName         string
	Restore          bool
	Key              string
	CryptoPrivateKey string
	ConfigFile       string
	TrustedSubnet    string
	GAddress         string
	TLSCertFile      string
	TLSKeyFile       string
}

type GRPCConfig struct {
	GAddress    string `env:"G_ADDRESS" json:"g_address"`
	TLSCertFile string `env:"CERT_FILE" json:"cert_file"`
	TLSKeyFile  string `env:"KEY_FILE" json:"key_file"`
}

// UseServerStartParams - assigning configuration environment variables.
func UseServerStartParams() ServerConfig {
	var c ServerConfig

	flag.StringVar(&c.DB, "d", "", "connection string to postgres db")
	flag.StringVar(&c.RunAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&c.LogLevel, "l", "info", "log level")
	flag.IntVar(&c.StoreInterval, "i", 300, "time interval in sec")
	flag.StringVar(&c.FileName, "f", "/tmp/metrics-db.json", "full file path")
	flag.BoolVar(&c.Restore, "r", true, "is load previously saved values")
	flag.StringVar(&c.Key, "k", "", "key")
	flag.StringVar(&c.CryptoPrivateKey, "crypto-key", "", "path to RSA private key file in PEM format")
	flag.StringVar(&c.ConfigFile, "c", "", "path to configuration file")
	flag.StringVar(&c.TrustedSubnet, "t", "", "trusted subnet")
	flag.StringVar(&c.GAddress, "g", "localhost:8081", "grpc server address")
	flag.StringVar(&c.TLSCertFile, "tls-cert", "", "tls cert file")
	flag.StringVar(&c.TLSKeyFile, "tls-key", "", "tls key file")

	flag.Parse()

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
	if envKey := os.Getenv("KEY"); envKey != "" {
		c.Key = envKey
	}
	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		c.CryptoPrivateKey = envCryptoKey
	}
	if envConfigFile := os.Getenv("CONFIG"); envConfigFile != "" {
		c.ConfigFile = envConfigFile
	}
	if envTrustedSubnet := os.Getenv("TRUSTED_SUBNET"); envTrustedSubnet != "" {
		c.TrustedSubnet = envTrustedSubnet
	}

	if envGAddress := os.Getenv("G_ADDRESS"); envGAddress != "" {
		c.GAddress = envGAddress
	}
	if envTLSCertFile := os.Getenv("CERT_FILE"); envTLSCertFile != "" {
		c.TLSCertFile = envTLSCertFile
	}
	if envTLSKeyFile := os.Getenv("KEY_FILE"); envTLSKeyFile != "" {
		c.TLSKeyFile = envTLSKeyFile
	}

	return c
}
