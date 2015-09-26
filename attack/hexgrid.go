package attack

import "math"

/*
 ========== HEX GRID ==========

       0, 1
 -1, 0       1, 1
       0, 0
 -1,-1       1, 0
       0,-1

circ = radius*6

*/
var HEXDIRS = map[int][2]int{
	0: [2]int{1, 0},
	1: [2]int{1, 1},
	2: [2]int{0, 1},
	3: [2]int{-1, 0},
	4: [2]int{-1, -1},
	5: [2]int{0, -1},
}

// Hex2Plane returns the x,y coords of the center of a hexagon
// chosen from a flat-side-up grid of "radius" sized hexagons,
// with the hexagon 0,0 having a center at 0,0
func Hex2Plane(radius int, coord [2]int) (plane [2]int) {
	x := 2 * radius * coord[0]
	// y^2 + (r/2)^2 = r^2
	// y = sqrt( r^2 - (r/2)^2 )
	y := float64(coord[1]*2) * math.Sqrt(float64(3*radius*radius)/4.0)
	return [2]int{x, int(y)}
}

// HexPolar returns the sector, radius, and 'sweep' of a grid
// coordinate: sectors are the galaxy split into pie pieces, sweep is how far radially into a sector a coordinate is
func HexPolar(grid [2]int) [3]int {
	if grid == [2]int{0, 0} {
		return [3]int{0, 0, 0}
	}
	r := HexDist(grid, [2]int{0, 0})
	circ := 6 * r
	for theta := 0; theta < circ; theta++ {
		if HexAngle2Grid(r, theta) == grid {
			sector := (theta / r) + 1
			sweep := theta % r
			return [3]int{sector, r, sweep}
		}
	}
	Log("Couldn't find HexPolar for", grid)
	return [3]int{-1, -1, -1}
}

// HexDist tells you how many steps from a to b
func HexDist(a, b [2]int) int {
	x := b[0] - a[0]
	y := b[1] - a[1]
	if x == 0 {
		if y >= 0 {
			return y
		} else {
			return -1 * y
		}
	} else if y == 0 {
		if x >= 0 {
			return x
		} else {
			return -1 * x
		}
	}
	if x < 0 && y < 0 {
		if x < y {
			return -1 * x
		} else {
			return -1 * y
		}
	} else if x > 0 && y > 0 {
		if x > y {
			return x
		} else {
			return y
		}
	}
	if x < 0 {
		x *= -1
	}
	if y < 0 {
		y *= -1
	}
	return y + x
}

// HexPath creates a path of 1 step grid coords, starting with A and ending with B
// HexPath uses HexPathSteps to calculate the least steps to take, and
// then uses StepSplit to interweave the steps as evenly as possible
func HexPath(a, b [2]int) [][2]int {
	d := HexDist(a, b)
	if d == 0 {
		return [][2]int{a}
	} else if d == 1 {
		return [][2]int{a, b}
	}
	stepper := a
	steps := HexPathSteps(a, b)
	var numLarge, numSmall int
	var sModX, sModY, lModX, lModY int
	if Abs(steps[0]) >= Abs(steps[1]) && Abs(steps[0]) >= Abs(steps[2]) {
		numLarge = Abs(steps[0])
		lModX, lModY = Dir(steps[0]), 0
		if steps[1] == 0 {
			numSmall = Abs(steps[2])
			sModX, sModY = Dir(steps[2]), Dir(steps[2])
		} else {
			numSmall = Abs(steps[1])
			sModX, sModY = 0, Dir(steps[1])
		}
	} else if Abs(steps[1]) >= Abs(steps[2]) {
		numLarge = Abs(steps[1])
		lModX, lModY = 0, Dir(steps[1])
		if steps[0] == 0 {
			numSmall = Abs(steps[2])
			sModX, sModY = Dir(steps[2]), Dir(steps[2])
		} else {
			numSmall = Abs(steps[0])
			sModX, sModY = Dir(steps[0]), 0
		}
	} else {
		numLarge = Abs(steps[2])
		lModX, lModY = Dir(steps[2]), Dir(steps[2])
		if steps[0] == 0 {
			numSmall = Abs(steps[1])
			sModX, sModY = 0, Dir(steps[1])
		} else {
			numSmall = Abs(steps[0])
			sModX, sModY = Dir(steps[0]), 0
		}
	}
	useSmall := StepSplit(numLarge, numSmall)
	r := [][2]int{a}
	for count := 0; count < d; count++ {
		var modX, modY int
		if useSmall[count] {
			modX, modY = sModX, sModY
		} else {
			modX, modY = lModX, lModY
		}
		stepper = [2]int{stepper[0] + modX, stepper[1] + modY}
		r = append(r, stepper)
	}
	return r
}

// HexPathSteps calculates the minimum steps from a to b
// HexPathSteps returns a [3]int with the count of:
// StepsRight(0,1) StepsUp(1,0) StepsUpRight(1,1)
// with negative numbers for backward steps
func HexPathSteps(a, b [2]int) (steps [3]int) {
	//steps right up upRight
	x := b[0] - a[0]
	y := b[1] - a[1]
	if x == 0 {
		steps[1] = y
	} else if y == 0 {
		steps[0] = x
	} else if x < 0 && y < 0 {
		if x < y {
			steps[2] = y
			steps[0] = x - y
		} else {
			steps[2] = x
			steps[1] = y - x
		}
	} else if x > 0 && y > 0 {
		if x > y {
			steps[2] = y
			steps[0] = x - y
		} else {
			steps[2] = x
			steps[1] = y - x
		}
	} else {
		steps[0] = x
		steps[1] = y
	}
	return
}

// StepSplit figures out how best to weave two types of steps
// It takes two ints representing counts of each steptype and is
// indifferent to what those steptypes are
// It returns a bool slice to use for determining which type to
// use in a sequence
func StepSplit(larger, smaller int) (useSmall []bool) {
	useSmall = make([]bool, larger+smaller)
	if smaller <= 0 {
		return
	}
	beat := (larger + smaller) / (smaller + 1)
	for i := 1; i < (smaller + 1); i++ {
		useSmall[i*beat] = true
	}
	return
}

func Abs(x int) int {
	if x > 0 {
		return x
	}
	return -1 * x
}

func Dir(x int) int {
	if x < 0 {
		return -1
	} else {
		return 1
	}
}

// HexAngle2Grid transforms polar step instructions (from the grid origin) into grid coordinates
// radius = steps outward directly rightward (1,0) from origin
// theta = steps counter-clockwise from this point, staying distance "radius" from the origin
func HexAngle2Grid(radius, theta int) (loc [2]int) {
	if radius == 0 {
		return [2]int{0, 0}
	}
	circ := 6 * radius
	for theta > circ {
		theta -= circ
	}
	for theta < 0 {
		theta += circ
	}
	// Which leg of the hexagon are we on?
	// (legs have length = radius)
	leg := theta / radius
	// How far along that leg are we?
	extra := theta % radius
	// Go out along the axis to the start of your leg
	extend := HEXDIRS[leg]
	extend = [2]int{extend[0] * radius, extend[1] * radius}
	// Go along the leg the appropriate distance
	legDir := HEXDIRS[(leg+2)%6]
	legDir = [2]int{legDir[0] * extra, legDir[1] * extra}
	// Viola!
	return [2]int{extend[0] + legDir[0], extend[1] + legDir[1]}
}
