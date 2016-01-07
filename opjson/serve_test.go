package opjson

import (
	"log"
	"net/http"
	"testing"
)

func TestFirst(t *testing.T) {
	log.Println("TEST FIRST")
}

func XTestSecond(t *testing.T) {
	const SERVPORT = ":8080"
	log.Println("STARTING SERVER AT", SERVPORT)
	http.HandleFunc("/serve", pageTestServeJson)
	http.HandleFunc("/read", pageTestReadJson)
	err := http.ListenAndServe(SERVPORT, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("STOPPING SERVER")
}
func pageTestServeJson(w http.ResponseWriter, r *http.Request) {
	g := &Game{}
	if !HttpServe(w, g) {
		log.Println("ERROR SERVING", g)
	}
	log.Println("SERVED", g)
}

func pageTestReadJson(w http.ResponseWriter, r *http.Request) {
	g := &Game{}
	if !HttpRead(w, r, g) {
		log.Println("ERROR READING", g)
	}
	log.Println("READ", g)
}
