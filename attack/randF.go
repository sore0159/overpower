package attack

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func shuffleWords(list []string) []string {
	if list == nil {
		return make([]string, 0)
	}
	l := len(list)
	if l < 2 {
		return list
	}
	newList := make([]string, l)
	newOrder := rand.Perm(l)
	for i, val := range newOrder {
		newList[val] = list[i]
	}
	return newList

}

func pick(n int) int {
	return rand.Intn(n) + 1
}

func shuffleShips(in []*Ship) (out []*Ship) {
	l := len(in)
	if l < 1 {
		return []*Ship{}
	} else if l == 1 {
		return []*Ship{in[0]}
	}
	out = make([]*Ship, l)
	for i, j := range rand.Perm(l) {
		out[i] = in[j]
	}
	return out
}
