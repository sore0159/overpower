package attack

import (
	"log"
	"testing"
)

func TestFirst(t *testing.T) {
	log.Println("TESTING")
}

func TestSecond(t *testing.T) {
	db, err := LoadDB()
	if err != nil {
		log.Println("Error getting db")
		return
	}
	_ = db
	log.Println("Got db")
	facs := map[string]string{
		"fac1": "user1",
		"fac2": "user2",
		"fac3": "user3",
	}
	gid, err := MakeGame(db, "Game Name", "user1", facs)
	if err != nil {
		log.Println("GAME 1 CONSTRUCTION FAILED")
		return
	}
	log.Println("GAME 1 MADE:", gid)
	gid2, err := MakeGame(db, "Game2", "user2", facs)
	if err != nil {
		log.Println("GAME 2 CONSTRUCTION FAILED")
		return
	}
	log.Println("GAME 2 MADE:", gid2)
	DelGame(db, gid)
	DelGame(db, gid2)
}
