package configs

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

func (c *ConfigAgent) GetPublicKey() (*rsa.PublicKey, error) {
	key, err := os.ReadFile(c.CryptoPublicKey)
	if err != nil {
		return nil, fmt.Errorf("Error read file : %s\n", err)
	}

	rsaKey, err := bytesToPublicKey(key)
	if err != nil {
		return nil, fmt.Errorf("Error bytes to public key : %s\n", err)
	}

	return rsaKey, nil
}

func bytesToPublicKey(b []byte) (*rsa.PublicKey, error) {
	var err error

	block, _ := pem.Decode(b)
	if block == nil {
		return nil, errors.New("unable to decode PEM private key")
	}

	blockBytes := block.Bytes
	public, err := x509.ParsePKCS1PublicKey(blockBytes)
	if err != nil {
		return nil, fmt.Errorf("Error parse PKCS1 public key: : %s\n", err)
	}

	return public, nil
}

func (c *ServerConfig) GetPrivateKey() (*rsa.PrivateKey, error) {
	key, err := os.ReadFile(c.CryptoPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Error read file : %s\n", err)
	}

	rsaKey, err := bytesToPrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("Error bytes to private key: %s\n", err)
	}

	return rsaKey, nil
}

func bytesToPrivateKey(b []byte) (*rsa.PrivateKey, error) {
	var err error

	block, _ := pem.Decode(b)
	if block == nil {
		return nil, errors.New("unable to decode PEM private key")
	}

	blockBytes := block.Bytes
	private, err := x509.ParsePKCS1PrivateKey(blockBytes)
	if err != nil {
		return nil, fmt.Errorf("Error parse PKCS1 private key: %s\n", err)
	}

	return private, nil
}
