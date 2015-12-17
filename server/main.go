package main

import (
	"database/sql"
	"log"
	"mule/overpower/db"
	"mule/users"
	"net/http"
)

const (
	SERVPORT = ":8080"
)

var (
	USERREG  *users.Registry
	ATTACKDB *sql.DB
)

func main() {
	var ok bool
	USERREG, ok = users.GetRegistry()
	if !ok {
		log.Fatal("Can't load userdata")
	}
	defer USERREG.Close()
	if ATTACKDB, ok = db.LoadDB(); !ok {
		log.Fatal("Can't load planet data")
	}
	defer ATTACKDB.Close()
	SetupMux()
	log.Println("STARTING SERVER AT", SERVPORT)
	err := http.ListenAndServe(SERVPORT, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("STOPPING SERVER")
}
