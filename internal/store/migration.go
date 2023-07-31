package store

func CreateTable(s *DBStorage) error {
	const query = `
	CREATE TABLE IF NOT EXISTS metrics (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
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
