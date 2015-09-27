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
	USERSPWS, err = LoadUserData()
	if err != nil {
		log.Fatal("Can't load userdata:", err)
	}
	TEMPGAME, err = Load(SAVEFL)
	if err != nil {
		log.Fatal("Can't load savefile:", err)
	}
	if TEMPGAME == nil {
		TEMPGAME = NewGame()
	}
	http.HandleFunc("/", gameMux)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("STATIC/"))))
	log.Println("STARTING SERVER AT", SERVPORT)
	err = http.ListenAndServe(SERVPORT, nil)
	if err != nil {
		Log(err)
	}
	log.Println("STOPPING SERVER")
}
