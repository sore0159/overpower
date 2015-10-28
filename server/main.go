package main

import (
	"log"
	"net/http"
)

const (
	DATADIR  = "DATA/"
	SERVPORT = ":8080"
)

func main() {
	var err error
	USERDB, err = LoadUserData()
	if err != nil {
		log.Fatal("Can't load userdata:", err)
	}
	defer USERDB.Close()
	ATTACKDB, err = LoadPlanetDB()
	if err != nil {
		log.Fatal("Can't load planet db:", err)
	}
	defer ATTACKDB.Close()
	http.HandleFunc("/", gameMux)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("STATIC/"))))
	log.Println("STARTING SERVER AT", SERVPORT)
	err = http.ListenAndServe(SERVPORT, nil)
	if err != nil {
		Log(err)
	}
	log.Println("STOPPING SERVER")
}
