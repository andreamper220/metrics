package storages

import "database/sql"

type DBStorage struct {
	*AbstractStorage
	Connection *sql.DB
}

func NewDBStorage(conn *sql.DB) *DBStorage {
	return &DBStorage{
		Connection: conn,
	}
}

func (dbs *DBStorage) WriteMetrics() error {
	return nil
}

func (dbs *DBStorage) ReadMetrics() error {
	return nil
}
