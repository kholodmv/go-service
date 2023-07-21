package store

import (
	"context"
	"github.com/kholodmv/go-service/internal/models"
)

type Storage interface {
	GetAllMetrics(ctx context.Context, size int64) []models.Metrics
	GetCountMetrics(ctx context.Context) int64
	GetValueMetric(ctx context.Context, typeM string, name string) (interface{}, bool)
	AddMetric(ctx context.Context, typeM string, value interface{}, name string) error

	GetAllMetricsJSON() []models.Metrics
	WriteAndSaveMetricsToFile(filename string) error
	RestoreFileWithMetrics(filename string)
	Ping() error
}
