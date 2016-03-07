package models

import (
	"database/sql"
	"errors"
	"mule/mydb/db"
)

var (
	ErrNoneFound = errors.New("none found")
	ErrNotUnique = errors.New("query did not result in unique match")
)

type DB struct {
	*sql.DB
}

func LoadDB() (*DB, error) {
	d, err := db.LoadDB(DB_USER, DB_PASS, PADB_NAME)
	if my, bad := Check(err, "loaddb  failure"); bad {
		return nil, my
	}
	return &DB{d}, nil
}

func (d *DB) NewManager() *Manager {
	return NewManager(d.DB)
}

func (d *DB) Transact(f func(*Manager) (error, error)) (logErr, failed error) {
	g := func(d db.DBer) error {
		m := NewManager(d)
		var revertE error
		logErr, revertE = f(m)
		if my, bad := Check(revertE, "manager transaction failure on execution"); bad {
			return my
		}
		revertE = m.Close()
		if my, bad := Check(revertE, "manager transaction failure on closure"); bad {
			return my
		}
		return nil
	}
	err := db.Transact(d.DB, g)
	if my, bad := Check(err, "managar transaction failed on db transact"); bad {
		return logErr, my
	}
	return logErr, nil
}
