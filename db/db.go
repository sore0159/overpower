package db

import (
	"database/sql"
	"mule/mydb"
)

type DB struct {
	db *sql.DB
}

func LoadDB() (DB, bool) {
	db, ok := mydb.LoadDB(DB_USER, DB_PASS, PADB_NAME)
	if !ok {
		return DB{nil}, false
	}
	return DB{db}, true
}

func (d DB) Update(list []mydb.Updater) bool {
	return mydb.Update(d.db, list)
}

func (d DB) Close() {
	d.db.Close()
}
