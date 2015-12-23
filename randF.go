package overpower

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
