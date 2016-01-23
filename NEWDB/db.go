package db

import (
	"database/sql"
	"errors"
	"mule/mybad"
	"mule/mydb"
)

var (
	Check        = mybad.BuildCheck("package", "overpower/db")
	ErrNoneFound = errors.New("none found")
	ErrNotUnique = errors.New("query did not result in unique match")
)

type DB struct {
	db *sql.DB
}

func LoadDB() (DB, error) {
	db, err := mydb.LoadDB(DB_USER, DB_PASS, PADB_NAME)
	if my, bad := Check(err, "loaddb  failure"); bad {
		return DB{nil}, my
	}
	return DB{db}, nil
}

func (d DB) Close() {
	d.db.Close()
}

// C is a convenience handler for conditions
// passed to the exporting functions
// Usage example:
//
// games, err := db.GetGames(C{"gid": 1})
type C []interface{}
