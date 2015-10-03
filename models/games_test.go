package models

import (
	"log"
	"testing"
)

func TestFirst(t *testing.T) {
	log.Println("TESTING")
}

func TestSecond(t *testing.T) {
	g := MakeGame("HI", "HI")
	if g != nil {
		log.Println("GAME MADE:", g)
		g.Destroy()
	}
}
