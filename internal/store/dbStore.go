package store

import (
	"context"
	"database/sql"
	"github.com/kholodmv/go-service/cmd/metrics"
	"github.com/kholodmv/go-service/internal/models"
	"log"
)

type DBStorage struct {
	db   *sql.DB
	path string
}

func NewStorage(path string) *DBStorage {
	s := &DBStorage{
		path: path,
	}
	db, err := sql.Open("postgres", s.path)
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	s.db = db
	s.createTable()

	return s
}

func (s *DBStorage) Ping() error {
	if err := s.db.Ping(); err != nil {
		log.Fatal("Failed to ping the database: ", err)
	}

	return nil
}

func (s *DBStorage) createTable() error {
	const query = `
	CREATE TABLE IF NOT EXISTS metrics (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		value DOUBLE PRECISION,
		delta BIGINT
	);
	`

	if _, err := s.db.Exec(query); err != nil {
		return err
	}
	return nil
}

func (s *DBStorage) GetAllMetrics(ctx context.Context, size int64) []models.Metrics {
	allM := make([]models.Metrics, 0, size)

	rows, err := s.db.QueryContext(ctx, "SELECT name, type, value, delta FROM metrics")
	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var m models.Metrics
		err = rows.Scan(&m.ID, &m.MType, &m.Value, &m.Delta)
		if err != nil {
			return nil
		}
		allM = append(allM, m)
	}

	if err = rows.Err(); err != nil {
		return nil
	}

	return allM
}

func (s *DBStorage) GetCountMetrics(ctx context.Context) int64 {
	row := s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) as count FROM metrics")
	var sumCount int64

	if err := row.Scan(&sumCount); err != nil {
		panic(err)
	}

	return sumCount
}

func (s *DBStorage) GetValueMetric(ctx context.Context, typeM string, name string) (interface{}, bool) {
	var row *sql.Row
	var mValue interface{}
	var ok = true
	if typeM == metrics.Gauge {
		row = s.db.QueryRowContext(ctx,
			"SELECT value FROM metrics WHERE name = $1 AND type = $2", name, typeM)

		if err := row.Scan(&mValue); err != nil {
			ok = false
		}
	}
	if typeM == metrics.Counter {
		row = s.db.QueryRowContext(ctx,
			"SELECT delta FROM metrics WHERE name = $1 AND type = $2", name, typeM)

		if err := row.Scan(&mValue); err != nil {
			ok = false
		}
	}

	return mValue, ok
}

func (s *DBStorage) AddMetric(ctx context.Context, typeM string, value interface{}, name string) error {
	if typeM == metrics.Gauge {
		result, err := s.db.ExecContext(ctx, "UPDATE metrics SET value = $2 WHERE name = $1", name, value)
		if err != nil {
			return err
		}

		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			if _, err = s.db.ExecContext(ctx, "INSERT INTO metrics (name, type, value) VALUES ($1, $2, $3)", name, typeM, value); err != nil {
				return err
			}
		}
	}
	if typeM == metrics.Counter {
		result, err := s.db.ExecContext(ctx, "SELECT name, type, value, delta FROM metrics WHERE name=$1", name)
		if err != nil {
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			if _, err = s.db.ExecContext(ctx, "INSERT INTO metrics (name, type, delta) VALUES ($1, $2, $3)", name, typeM, value); err != nil {
				return err
			}
			return nil
		}

		ms := models.Metrics{}
		row := s.db.QueryRowContext(ctx, "SELECT name, type, value, delta FROM metrics WHERE name=$1", name)
		if err = row.Scan(&ms.ID, &ms.MType, &ms.Value, &ms.Delta); err != nil {
			return err
		}

		v := value.(int64)
		count := v + *ms.Delta

		if _, err = s.db.ExecContext(ctx, "UPDATE metrics SET delta = $2 WHERE name = $1", name, count); err != nil {
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
