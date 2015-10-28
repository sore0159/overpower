package main

import (
	"log"
	"testing"
)

func TestFirst(t *testing.T) {
	log.Println("TESTING")
}

func TestSecond(t *testing.T) {
	USERSPWS, _ = LoadUserData()
	log.Println(USERSPWS)
	log.Println(ValidPW("one", "two"))
	log.Println(ValidPW("two", "two"))
}
