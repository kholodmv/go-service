package configs

import (
	"os"
	"testing"
)

func TestGetEnvironmentRunAddressVariable(t *testing.T) {
	var c ServerConfig
	c.RunAddress = "localhost:8080"
	err := os.Setenv("ADDRESS", c.RunAddress)
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	defer os.Unsetenv("ADDRESS") // Удаляем параметр окружения после теста.

	// Вызываем функцию для получения значения параметра окружения.
	result := UseServerStartParams()

	// Сравниваем полученное значение с ожидаемым.
	if result.RunAddress != c.RunAddress {
		t.Errorf("Expected value: %s, got: %s", c.RunAddress, result.RunAddress)
	}
}
