package attack

/*
 ========== HEX GRID ==========

       1, 0
  0,-1       1, 1
       0, 0
 -1,-1       0, 1
      -1, 0

circ = (r-1)*6


*/

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

func StepSplit(larger, smaller int) (useSmall []bool) {
	useSmall = make([]bool, larger+smaller)
	beat := (larger + smaller) / (smaller + 1)
	if (larger+smaller)%(smaller+1) != 0 {
		beat++
	}
	count := 0
	for i := 0; i < larger+smaller; i++ {
		if (i+1)%beat == 0 && count < smaller {
			count++
			useSmall[i] = true
		}
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

func HexAngle2Grid(radius, theta int) (loc [2]int) {
	if radius == 0 {
		return [2]int{0, 0}
	}
	leg := theta / radius // 0-5
	extra := theta % radius
	var mainX, mainY, boostX, boostY int
	switch leg {
	case 0: //  1, 0
		mainX, mainY, boostX, boostY = 1, 0, 0, -1
	case 1: //  0,-1
		mainX, mainY, boostX, boostY = 0, -1, -1, -1
	case 2: // -1,-1
		mainX, mainY, boostX, boostY = -1, -1, -1, 0
	case 3: // -1, 0
		mainX, mainY, boostX, boostY = -1, 0, 0, 1
	case 4: //  0, 1
		mainX, mainY, boostX, boostY = 0, 1, 1, 1
	case 5: //  1, 1
		mainX, mainY, boostX, boostY = 1, 1, 1, 0
	}
	loc = [2]int{mainX*radius + boostX*extra, mainY*radius + boostY*extra}
	return
}
