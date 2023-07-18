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
	err := s.initDB()
	if err != nil {
		log.Printf("Database initialization error: %v", err)
	}
	return s
}

func (s *Storage) initDB() error {
	db, err := sql.Open("postgres", s.storagePath)
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	s.db = db

	return nil
}

func (s *Storage) Ping() error {
	err := s.db.Ping()
	if err != nil {
		log.Fatal("Failed to ping the database: ", err)
	}

	return nil
}
