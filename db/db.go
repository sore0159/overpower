package db

import (
	"database/sql"
	"mule/mydb"
)

func LoadDB() (*sql.DB, bool) {
	return mydb.LoadDB(DB_USER, DB_PASS, PADB_NAME)
}
