package models

import (
	"log"
	"testing"
)

func TestFirst(t *testing.T) {
	log.Println("TESTING")
}

func TestSecond(t *testing.T) {
	facs := map[string]string{
		"fac1": "user1",
		"fac2": "user2",
		"fac3": "user3",
	}
	g := MakeGame("Game Name", "user1", facs)
	if g == nil {
		log.Println("GAME 1 CONSTRUCTION FAILED")
		return
	}
	log.Println("GAME 1 MADE:", g)
	g2 := MakeGame("Game2", "user2", facs)
	if g2 == nil {
		log.Println("GAME 2 CONSTRUCTION FAILED")
		return
	}
	log.Println("GAME 2 MADE:", g2)
	DelGame(g.GID)
	DelGame(g2.GID)
}
