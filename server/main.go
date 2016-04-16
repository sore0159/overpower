package main

import (
	"math/rand"
	"mule/overpower/models"
	"mule/users"
	"net/http"
	"time"
)

const (
	DATADIR  = "DATA/"
	SERVPORT = ":8080"
)

var (
	USERREG *users.Registry
	OPDB    *models.DB
	DBLOCK  bool
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var err error
	USERREG, err = users.GetRegistry()
	if my, bad := Check(err, "failed to get registry"); bad {
		Log(my)
		panic(my)
	}
	defer USERREG.Close()
	OPDB, err = models.LoadDB()
	if my, bad := Check(err, "failed to get op database"); bad {
		Log(my)
		panic(my)
	}
	defer OPDB.Close()
	go AutoTimer()
	SetupMux()
	Announce("STARTING SERVER AT", SERVPORT)
	err = http.ListenAndServe(SERVPORT, nil)
	if my, bad := Check(err, "server failure"); bad {
		Log(my)
	}
	Announce("STOPPING SERVER")
}
