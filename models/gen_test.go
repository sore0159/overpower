package models

import (
	"log"
	//	"mule/hexagon"
	//	"mule/mybad"
	//	"mule/overpower"
	"testing"
)

func TestOne(t *testing.T) {
	log.Println("TEST ONE")
}

func TestTwo(t *testing.T) {
	db, err := LoadDB()
	ErrCheck(err)
	log.Println("Loaded DB")
	logE, failE := db.Transact(MakeTest)
	ErrCheck(failE)
	ErrCheck(logE)
	log.Println("MakeTest complete!")
}

func MakeTest(m *Manager) (logE, failE error) {
	return nil, nil
}
