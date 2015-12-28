package db

import (
	"database/sql"
	"fmt"
	"mule/mydb"
	"testing"
)

const UPDATETABLES = true

func MakeTables(db *sql.DB) (ok bool) {
	queries := []string{}
	queries = append(queries, `create table games(
	gid SERIAL PRIMARY KEY,
	owner varchar(20) NOT NULL UNIQUE,
	name varchar(20) NOT NULL,
	turn int NOT NULL,
	password varchar(20)
);`)
	queries = append(queries, `create table factions(
	gid integer NOT NULL REFERENCES games ON DELETE CASCADE,
	fid SERIAL NOT NULL,
	owner varchar(20) NOT NULL,
	name varchar(20) NOT NULL,
	done bool NOT NULL,
	UNIQUE(gid, owner),
	PRIMARY KEY(gid, fid)
);`)
	queries = append(queries, `create table mapviews(
	gid integer NOT NULL REFERENCES games ON DELETE CASCADE,
	fid integer NOT NULL,
	center point NOT NULL,
	zoom int NOT NULL,
	focus point,
	FOREIGN KEY(gid, fid) REFERENCES factions ON DELETE CASCADE,
	PRIMARY KEY (gid, fid)
);`)
	queries = append(queries, `create table planets(
	gid integer NOT NULL REFERENCES games ON DELETE CASCADE,
	pid integer NOT NULL,
	name varchar(20) NOT NULL,
	loc point NOT NULL,
	controller int,
	inhabitants int NOT NULL,
	resources int NOT NULL,
	parts int NOT NULL,
	UNIQUE(gid, name),
	PRIMARY KEY(gid, pid),
	FOREIGN KEY(gid, controller) REFERENCES factions ON DELETE CASCADE
);`)
	queries = append(queries, `create table planetviews(
	gid integer NOT NULL REFERENCES games ON DELETE CASCADE,
	fid integer NOT NULL,
	pid integer NOT NULL,
	name varchar(20) NOT NULL,
	loc point NOT NULL,
	turn int NOT NULL,
	controller int,
	inhabitants int NOT NULL,
	resources int NOT NULL,
	parts int NOT NULL,
	FOREIGN KEY(gid, fid) REFERENCES factions ON DELETE CASCADE,
	FOREIGN KEY(gid, controller) REFERENCES factions ON DELETE CASCADE,
	FOREIGN KEY(gid, pid) REFERENCES planets ON DELETE CASCADE,
	PRIMARY KEY(gid, fid, pid)
);`)
	queries = append(queries, `create table orders(
	gid integer NOT NULL REFERENCES games ON DELETE CASCADE,
	fid integer NOT NULL,
	source integer NOT NULL,
	target integer NOT NULL,
	size integer NOT NULL,
	FOREIGN KEY(gid, fid) REFERENCES factions ON DELETE CASCADE,
	FOREIGN KEY(gid, source) REFERENCES planets ON DELETE CASCADE,
	FOREIGN KEY(gid, target) REFERENCES planets ON DELETE CASCADE,
	PRIMARY KEY(gid, fid, source, target)
);`)
	queries = append(queries, `create table ships(
	gid integer NOT NULL REFERENCES games ON DELETE CASCADE,
	fid int NOT NULL,
	sid SERIAL NOT NULL,
	size int NOT NULL,
	launched int NOT NULL,
	path point[] NOT NULL,
	FOREIGN KEY(gid, fid) REFERENCES factions ON DELETE CASCADE,
	PRIMARY KEY(gid, fid, sid)
);`)
	queries = append(queries, `create table shipviews(
	gid integer NOT NULL REFERENCES games ON DELETE CASCADE,
	fid integer NOT NULL,
	controller integer NOT NULL,
	sid integer NOT NULL,
	turn integer NOT NULL,
	loc point,
	trail point[],
	size int NOT NULL,
	FOREIGN KEY(gid, controller) REFERENCES factions ON DELETE CASCADE,
	FOREIGN KEY(gid, fid) REFERENCES factions ON DELETE CASCADE,
	PRIMARY KEY(gid, fid, turn, sid)
);`)
	for i, query := range queries {
		if !mydb.ExecIf(db, query) {
			fmt.Println("Failed table creation", i)
			return false
		}
		fmt.Println("Table update", i, "passed")
	}
	return true
}

func DropTables(db *sql.DB) (ok bool) {
	tables := "games, planets, factions, mapviews, ships, shipviews, planetviews, orders"
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", tables)
	return mydb.ExecIf(db, query)
}

func TestUpdateTables(t *testing.T) {
	if UPDATETABLES {
		d, ok := LoadDB()
		if !ok {
			fmt.Println("FAILED TO LOAD DB")
			return
		}
		db := d.db
		if DropTables(db) {
			fmt.Println("Dropped tables!")
		} else {
			fmt.Println("failed dropped tables!")
		}
		if MakeTables(db) {
			fmt.Println("Made tables!")
		} else {
			fmt.Println("failed make tables!")
		}
	}
}
