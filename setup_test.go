package overpower

import (
	"log"
	"testing"
)

func TestFirst(t *testing.T) {
	log.Println("TEST ONE")
}

func TestSecond(t *testing.T) {
	fids := []int{1, 2, 3, 4, 5, 6, 7, 8}
	planets := MakeGalaxy(fids)
	log.Println("FIDS:", fids)
	log.Println("MADE PLANETS:", planets)
}
