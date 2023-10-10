// Package store - store.go - contains interface with methods for work with db and memory.
package store

import (
	"context"

	"github.com/kholodmv/go-service/internal/models"
)

// Storage interface with methods for work with db and memory.
type Storage interface {
	// GetAllMetrics the function tries to get all the metrics.
	GetAllMetrics(ctx context.Context, size int64) ([]models.Metrics, error)
	// GetCountMetrics the function tries to get count of all metrics.
	GetCountMetrics(ctx context.Context) (int64, error)
	// GetValueMetric the function tries to get value of concrete metric by type and name.
	GetValueMetric(ctx context.Context, typeM string, name string) (interface{}, error)
	// AddMetric the function tries to add value of concrete metric.
	AddMetric(ctx context.Context, typeM string, value interface{}, name string) error

	// GetAllMetricsJSON the function tries to get all the metrics in json format.
	GetAllMetricsJSON() []models.Metrics
	// WriteAndSaveMetricsToFile the function writes metrics to a file and saves them there.
	WriteAndSaveMetricsToFile(filename string) error
	// RestoreFileWithMetrics the function opens a file and saves metrics to it.
	RestoreFileWithMetrics(filename string)
	// Ping to db.
	Ping() error
}
