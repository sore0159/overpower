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
	realDB *sql.DB
	tx     *sql.Tx
}

func LoadDB() (DB, error) {
	db, err := mydb.LoadDB(DB_USER, DB_PASS, PADB_NAME)
	if my, bad := Check(err, "loaddb  failure"); bad {
		return DB{nil, nil}, my
	}
	return DB{db, nil}, nil
}

func (d DB) Close() error {
	if d.realDB == nil {
		return errors.New("cannot close: in transaction")
	}
	return d.realDB.Close()
}

func (d DB) Begin() (DB, error) {
	if d.realDB == nil {
		return DB{nil, nil}, errors.New("cannot begin: already in transaction")
	}
	tx, err := d.realDB.Begin()
	if my, bad := Check(err, "op/db begin failure"); bad {
		return DB{nil, nil}, my
	}
	return DB{nil, tx}, nil
}
func (d DB) InTrans() bool {
	return d.realDB == nil
}

func (d DB) db() mydb.DBer {
	if d.tx == nil {
		return d.realDB
	}
	return d.tx
}

func (d DB) Commit() error {
	if d.tx == nil {
		return errors.New("cannot commit: not in transaction")
	}
	return d.tx.Commit()
}
func (d DB) Rollback() error {
	if d.tx == nil {
		return errors.New("cannot rollback: not in transaction")
	}
	return d.tx.Rollback()
}

func (d DB) Transact(f func(DB) error) error {
	tx, err := d.Begin()
	if my, bad := Check(err, "transaction begin failure"); bad {
		return my
	}
	err = f(tx)
	if my, bad := Check(err, "transact execute failure"); bad {
		if err2 := tx.Rollback(); err2 != nil {
			my.AddContext("rollback failure", err2)
		}
		return my
	} else {
		if my, bad := Check(tx.Commit(), "transact commit failure"); bad {
			return my
		}
		return nil
	}
}

// C is a convenience handler for conditions
// passed to the exporting functions
// Usage example:
//
// games, err := db.GetGames(C{"gid": 1})
type C []interface{}

func (d DB) mustExec(query string, args ...interface{}) error {
	return mydb.ExecCheck(d.db().Exec(query, args...))
}

func IsUnique(x int) error {
	switch x {
	case 0:
		return ErrNoneFound
	case 1:
		return nil
	default:
		return ErrNotUnique
	}
}
