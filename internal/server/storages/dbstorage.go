package storages

import "database/sql"

type DBStorage struct {
	conn *sql.DB
}
