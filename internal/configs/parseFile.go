package configs

import (
	"bytes"
	"encoding/json"
	"os"
)

func (c *ConfigAgent) ParseFile(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(bytes.NewReader(data)).Decode(c); err != nil {
		return err
	}

	return nil
}

func (c *ServerConfig) ParseFile(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(bytes.NewReader(data)).Decode(c); err != nil {
		return err
	}

	return nil
}
