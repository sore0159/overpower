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
	g := attack.MakeGame([]string{"test"})
	MakeMap(g.Factions[1])
}
