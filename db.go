package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const PWFILE = DATADIR + "passwords.txt"

var USERDB *sql.DB

func LoadUserData() (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASS, USERDB_NAME))
	if err != nil {
		return nil, Log(err)
	}
	err = db.Ping()
	if err != nil {
		return nil, Log(err)
	}
	return db, nil
}
