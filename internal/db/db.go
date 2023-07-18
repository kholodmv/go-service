package db

import (
	"database/sql"
	"log"
)

type Storage struct {
	db          *sql.DB
	storagePath string
}

type DBStorage interface {
	Ping() error
}

func NewStorage(path string) DBStorage {
	s := &Storage{
		storagePath: path,
	}
	db, err := sql.Open("postgres", s.storagePath)
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	s.db = db
	s.createTable()

	return s
}

func (s *Storage) Ping() error {
	err := s.db.Ping()
	if err != nil {
		log.Fatal("Failed to ping the database: ", err)
	}

	return nil
}

func (s *Storage) createTable() error {
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

/*func (s *Storage) GetAllMetrics() error {
	row := s.db.QueryRowContext(context.Background(),
		"SELECT name, type, value, delta "+
			"FROM videos ORDER BY likes DESC LIMIT 1")
	var (
		title  string
		likes  int
		comdis bool
	)
	// порядок переменных должен соответствовать порядку колонок в запросе
	err = row.Scan(&title, &likes, &comdis)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s | %d | %t \r\n", title, likes, comdis)
}
*/
