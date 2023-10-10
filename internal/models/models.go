// Package models contains information about models.
package models

// Metrics struct include information about metrics.
type Metrics struct {
	ID    string   `json:"id"`              // metric name
	MType string   `json:"type"`            // parameter taking the value gauge or counter
	Delta *int64   `json:"delta,omitempty"` // metric value in case of transfer counter
	Value *float64 `json:"value,omitempty"` // metric value in case of transmitting gauge
}
