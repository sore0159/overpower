package attack

import (
	"fmt"
	"testing"
)

func TestFirst(t *testing.T) {
	fmt.Println("TESTING")
}

func TestMake(t *testing.T) {
	g := MakeGame([]string{"1"})
	g = MakeGame([]string{"1", "2"})
	fmt.Println(g)
}
func XTestOrder(t *testing.T) {
	f := NewFaction()
	f.AddOrder(0, 0, 0, 0, 0)
	fmt.Println(f.BuildOrders)
	f.DropOrder(0)
	fmt.Println(f.BuildOrders)
}

func XTestDist(t *testing.T) {
	for _, a := range [][2]int{[2]int{0, 0}, [2]int{3, 1}, [2]int{-10, 60}, [2]int{-40, -50}, [2]int{5, -4}, [2]int{5, 5}} {
		for _, b := range [][2]int{[2]int{2, 5}, [2]int{4, 1}, [2]int{-2, 5}, [2]int{-2, -5}, [2]int{2, -5}} {
			fmt.Println("A:", a, "B:", b, "Dist:", HexDist(a, b))
			p := HexPath(a, b)
			fmt.Println("A:", a, "B:", b, "Path:", p)
		}
	}
}

func XTestSplitter(t *testing.T) {
	fmt.Println("UseSmall:", StepSplit(8, 2))
	fmt.Println("UseSmall:", StepSplit(4, 2))
}
