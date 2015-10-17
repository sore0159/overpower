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
	g := MakeGame(DB, "Game Name", "user1", facs)
	if g == nil {
		log.Println("GAME 1 CONSTRUCTION FAILED")
		return
	}
	log.Println("GAME 1 MADE:", g, g.Owner())
	g2 := MakeGame(DB, "Game2", "user2", facs)
	if g2 == nil {
		log.Println("GAME 2 CONSTRUCTION FAILED")
		return
	}
	log.Println("GAME 2 MADE:", g2, g2.Owner())
	DelGame(DB, g.ID())
	DelGame(DB, g2.ID())
}

func TestThird(t *testing.T) {
	facs := map[string]string{
		"fac1": "user1",
		"fac2": "user2",
		"fac3": "user3",
	}
	g := MakeGame(DB, "Game Name", "user1", facs)
	if g == nil {
		log.Println("GAME 1 CONSTRUCTION FAILED")
		return
	}
	log.Println("GAME 1 MADE:", g, g.Owner())
	for fid, f := range g.Factions() {
		log.Println("Faction", fid, ":", f.Name())
	}
	DelGame(DB, g.ID())
}
