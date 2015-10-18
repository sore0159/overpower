package attack

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	_ "github.com/lib/pq"
	"mule/mylog"
	"strconv"
	"strings"
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

type Point [2]int

func (p *Point) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return Log("Bad value scanned to point:", value)
	}
	parts := strings.Split(string(bytes), ",")
	x, err := strconv.Atoi(parts[0][1:])
	if err != nil {
		return Log("Bad point scan value", parts, err)
	}
	y, err := strconv.Atoi(parts[1][:len(parts[1])-1])
	if err != nil {
		return Log("Bad point scan value", parts, err)
	}
	p[0], p[1] = x, y
	return nil
}

func (p *Point) Value() (driver.Value, error) {
	return (*p).String(), nil
}

func (p Point) String() string {
	return fmt.Sprintf("POINT(%d,%d)", p[0], p[1])
}
