package planetattack

import (
	"database/sql"
	//	"database/sql/driver"
	"fmt"
	_ "github.com/lib/pq"
	"mule/mylog"
)

var (
	Log = mylog.Err
)

func init() {
	mylog.InitDefaults()
}

func LoadDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASS, PADB_NAME))
	if err != nil {
		return nil, Log(err)
	}
	err = db.Ping()
	if err != nil {
		return nil, Log(err)
	}
	return db, nil
}
