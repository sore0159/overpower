package main

import (
	"log"
	"mule/overpower/db"
	"mule/users"
	"net/http"
)

const (
	SERVPORT = ":8080"
)

var (
	USERREG *users.Registry
	OPDB    db.DB
	DBLOCK  bool
)

func main() {
	var ok bool
	USERREG, ok = users.GetRegistry()
	if !ok {
		log.Fatal("Can't load userdata")
	}
	defer USERREG.Close()
	if OPDB, ok = db.LoadDB(); !ok {
		log.Fatal("Can't load planet data")
	}
	defer OPDB.Close()
	go AutoTimer()
	SetupMux()
	log.Println("STARTING SERVER AT", SERVPORT)
	err := http.ListenAndServe(SERVPORT, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("STOPPING SERVER")
}
