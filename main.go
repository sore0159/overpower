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
	http.HandleFunc("/", gameMux)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("STATIC/"))))
	log.Println("STARTING SERVER AT", SERVPORT)
	http.ListenAndServe(SERVPORT, nil)
}
