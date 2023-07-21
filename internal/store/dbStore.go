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
	err := s.db.Ping()
	if err != nil {
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

	_, err := s.db.Exec(query)
	return err
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

	err = rows.Err()
	if err != nil {
		return nil
	}
	return allM
}

func (s *DBStorage) GetCountMetrics(ctx context.Context) int64 {
	row := s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) as count FROM metrics")
	var count int64
	err := row.Scan(&count)
	if err != nil {
		panic(err)
	}
	return count
}

func (s *DBStorage) GetValueMetric(ctx context.Context, typeM string, name string) (interface{}, bool) {
	var row *sql.Row
	var mValue interface{}
	var ok = true
	if typeM == metrics.Gauge {
		row = s.db.QueryRowContext(ctx,
			"SELECT value FROM metrics WHERE name = $1 AND type = $2", name, typeM)
		err := row.Scan(&mValue)
		if err != nil {
			ok = false
		}
	}
	if typeM == metrics.Counter {
		row = s.db.QueryRowContext(ctx,
			"SELECT delta FROM metrics WHERE name = $1 AND type = $2", name, typeM)
		err := row.Scan(&mValue)
		if err != nil {
			ok = false
		}
	}

	return mValue, ok
}

func (s *DBStorage) UpdateMetric(ctx context.Context, typeM string, value interface{}, name string) error {
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
			_, err = s.db.ExecContext(ctx, "INSERT INTO metrics (name, type, value) VALUES ($1, $2, $3)", name, typeM, value)
			if err != nil {
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
			_, err = s.db.ExecContext(ctx, "INSERT INTO metrics (name, type, delta) VALUES ($1, $2, $3)", name, typeM, value)
			if err != nil {
				return err
			}
			return nil
		}

		ms := models.Metrics{}
		row := s.db.QueryRowContext(ctx, "SELECT name, type, value, delta FROM metrics WHERE name=$1", name)
		err = row.Scan(&ms.ID, &ms.MType, &ms.Value, &ms.Delta)
		if err != nil {
			return err
		}

		v := value.(int64)
		count := v + *ms.Delta

		_, err = s.db.ExecContext(ctx, "UPDATE metrics SET delta = $2 WHERE name = $1", name, count)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *DBStorage) AddMetric(ctx context.Context, typeM string, value interface{}, name string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if typeM == metrics.Gauge {
		_, err = s.db.ExecContext(ctx,
			"INSERT INTO metrics (name, type, value)"+
				" VALUES($1,$2,$3)", name, typeM, value)
		if err != nil {
			return err
		}
	}
	if typeM == metrics.Counter {
		_, err = s.db.ExecContext(ctx,
			"INSERT INTO metrics (name, type, delta)"+
				" VALUES($1,$2,$3)", name, typeM, value)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *DBStorage) GetAllMetricsJSON() []models.Metrics {
	return nil
}

func (s *DBStorage) WriteAndSaveMetricsToFile(filename string) error {
	return nil
}

func (s *DBStorage) RestoreFileWithMetrics(filename string) {

}
