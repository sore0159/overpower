package attack

import (
	"fmt"
	"testing"
)

func TestFirst(t *testing.T) {
	fmt.Println("TESTING")
}

func XTestHex2Plane(t *testing.T) {
	for _, a := range [][2]int{[2]int{0, 0}, [2]int{3, 1}, [2]int{-10, 60}, [2]int{-40, -50}, [2]int{5, -4}, [2]int{5, 5}} {
		fmt.Println("Point:", a, "to plane:", Hex2Plane(30, a))
	}
	for _, b := range [][2]int{[2]int{2, 5}, [2]int{4, 1}, [2]int{-2, 5}, [2]int{-2, -5}, [2]int{2, -5}} {
		fmt.Println("Point:", b, "to plane:", Hex2Plane(30, b))
	}
}

func TestPath(t *testing.T) {
	fmt.Println("PATH:", HexPath([2]int{19, 32}, [2]int{0, 0}))
}

func XTestMake(t *testing.T) {
	g := MakeGame([]string{"1"})
	g = MakeGame([]string{"1", "2"})
	fmt.Println(g)
}

func TestDist(t *testing.T) {
	for _, a := range [][2]int{[2]int{0, 0}, [2]int{3, 1}, [2]int{-10, 60}, [2]int{-40, -50}, [2]int{5, -4}, [2]int{5, 5}} {
		for _, b := range [][2]int{[2]int{2, 5}, [2]int{4, 1}, [2]int{-2, 5}, [2]int{-2, -5}, [2]int{2, -5}} {
			//		fmt.Println("A:", a, "B:", b, "Dist:", HexDist(a, b))
			fmt.Println("------------------------------------------")
			d := HexDist(a, b)
			p := HexPath(a, b)

			fmt.Println("A:", a, "B:", b, "PathD:", len(p), "Path:", p)
			if p[len(p)-1] != b || len(p) != d+1 {
				fmt.Println("============= FAILED =============")
			}
		}
	}
}

func XTestHexAngle(t *testing.T) {
	for _, rT := range [][2]int{[2]int{3, 4}, [2]int{6, 40}, [2]int{5, 80}, [2]int{100, 300}, [2]int{58, 234}} {
		fmt.Println("R, T:", rT)
		fmt.Println("X, Y:", HexAngle2Grid(rT[0], rT[1]))
	}
}

func XTestSplitter(t *testing.T) {
	fmt.Println("UseSmall:", StepSplit(8, 2))
	fmt.Println("UseSmall:", StepSplit(4, 2))
}
