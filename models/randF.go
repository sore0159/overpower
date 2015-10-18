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

func RandBigPlLoc(bigRange int) [2]int {
	loc := [2]int{pick(bigRange), pick(bigRange)}
	return loc
}

func RandLittlePlLoc(bigRange, littleRange int) [2]int {
	loc := [2]int{bigRange + pick(littleRange-bigRange), bigRange + pick(littleRange-bigRange)}
	return loc
}

func RandHomePlLoc(bigRange, littleRange, i, numF int) [2]int {
	loc := [2]int{bigRange + pick(littleRange-bigRange), bigRange + pick(littleRange-bigRange)}
	return loc
}
