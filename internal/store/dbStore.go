// Package store - dbStore.go - contains the logic for working with the database.
package store

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	"github.com/kholodmv/go-service/cmd/metrics"
	"github.com/kholodmv/go-service/internal/models"
)

// DBStorage struct include params database and logger.
type DBStorage struct {
	db  *sql.DB
	log *zap.SugaredLogger
}

// NewStorage function which initialize db storage.
func NewStorage(db *sql.DB, log *zap.SugaredLogger) *DBStorage {
	s := &DBStorage{
		db:  db,
		log: log,
	}

	return s
}

// Ping the function checks whether there is a connection to the database.
func (s *DBStorage) Ping() error {
	if err := s.db.Ping(); err != nil {
		s.log.Errorf("Failed to ping the database: %v", err)
		return err
	}

	return nil
}

// GetAllMetrics the function tries to get all the metrics contained in the 'metrics' table.
func (s *DBStorage) GetAllMetrics(ctx context.Context, size int64) ([]models.Metrics, error) {
	allM := make([]models.Metrics, 0, size)

	rows, err := s.db.QueryContext(ctx, "SELECT name, type, value, delta FROM metrics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m models.Metrics
		err = rows.Scan(&m.ID, &m.MType, &m.Value, &m.Delta)
		if err != nil {
			s.log.Errorf("Error copy column values into variables: %v", err)
			return nil, err
		}
		allM = append(allM, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return allM, nil
}

// GetCountMetrics the function tries to get count of all metrics contained in the 'metrics' table.
func (s *DBStorage) GetCountMetrics(ctx context.Context) (int64, error) {
	row := s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) as count FROM metrics")
	var sumCount int64

	if err := row.Scan(&sumCount); err != nil {
		s.log.Errorf("Error copy column values into variables: %v", err)
		return sumCount, err
	}

	return sumCount, nil
}

// GetValueMetric the function tries to get value of concrete metric by type and name contained in the 'metrics' table.
func (s *DBStorage) GetValueMetric(ctx context.Context, typeM string, name string) (interface{}, error) {
	var row *sql.Row
	var mValue interface{}
	if typeM == metrics.Gauge {
		row = s.db.QueryRowContext(ctx,
			"SELECT value FROM metrics WHERE name = $1 AND type = $2", name, typeM)

		if err := row.Scan(&mValue); err != nil {
			s.log.Errorf("Error copy column values into variables.: %v", err)
			return mValue, err
		}
	}
	if typeM == metrics.Counter {
		row = s.db.QueryRowContext(ctx,
			"SELECT delta FROM metrics WHERE name = $1 AND type = $2", name, typeM)

		if err := row.Scan(&mValue); err != nil {
			s.log.Errorf("Error copy column values into variables.: %v", err)
			return mValue, err
		}
	}

	return mValue, nil
}

// AddMetric the function tries to add value of concrete metric to 'metrics' table.
func (s *DBStorage) AddMetric(ctx context.Context, typeM string, value interface{}, name string) error {
	if typeM == metrics.Gauge {
		_, err := s.db.ExecContext(ctx, "INSERT INTO metrics (name, type, value) VALUES ($1, $2, $3) ON CONFLICT (name) DO UPDATE SET value = $3", name, typeM, value)
		if err != nil {
			return err
		}
	}
	if typeM == metrics.Counter {
		_, err := s.db.ExecContext(ctx, "INSERT INTO metrics (name, type, delta) VALUES ($1, $2, $3) ON CONFLICT (name) DO UPDATE SET delta = metrics.delta + $3;", name, typeM, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *DBStorage) GetAllMetricsJSON() []models.Metrics {
	return nil
}

func (s *DBStorage) WriteAndSaveMetricsToFile(filename string) error {
	return nil
}

func (s *DBStorage) RestoreFileWithMetrics(filename string) {

}
