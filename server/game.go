package main

import (
	"mule/planetattack"
	"unicode"
)

type Game struct {
	planetattack.Game
}

func ValidGameName(gName string) bool {
	if gName == "" || len(gName) > 15 {
		return false
	}
	for _, rn := range gName {
		if !unicode.In(rn, unicode.L, unicode.N) {
			return false
		}
	}
	return true
}
