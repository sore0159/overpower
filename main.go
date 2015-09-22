package main

import (
	"log"
	"mule/mylog"
	"net/http"
)

const DATADIR = "DATA/"

var Log = mylog.Err

func init() {
	mylog.SetErr(DATADIR + "errors.txt")
}

func main() {
	http.HandleFunc("/", gameMux)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("STATIC/"))))
	log.Println("STARTING SERVER")
	http.ListenAndServe(":8080", nil)
}
