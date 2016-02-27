package overpower

import (
	"math/rand"
)

func shuffleInts(list []int) []int {
	if list == nil {
		return make([]int, 0)
	}
	l := len(list)
	if l < 2 {
		return list
	}
	newList := make([]int, l)
	newOrder := rand.Perm(l)
	for i, val := range newOrder {
		newList[val] = list[i]
	}
	return newList
}

func pick(n int) int {
	return rand.Intn(n) + 1
}

func coin() bool {
	return rand.Intn(2) == 0
}

func RandF(x float64) float64 {
	return rand.Float64() * x
}
