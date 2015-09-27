package mapping

import (
	"fmt"
	"mule/planetattack/attack"
	"testing"
)

func TestFirst(t *testing.T) {
	fmt.Println("TESTING")
}

func TestSecond(t *testing.T) {
	n := 2
	g := attack.MakeGame([]string{"test", "test6", "test7", "test8", "test5", "test2", "test3", "test4"}[:n])
	fmt.Println("MAPPING")
	MakeMap(g.Factions[1])
}
