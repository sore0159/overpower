package planetattack

import (
	"mule/hexagon"
	"strings"
)

func (g *Game) MakeGalaxy(fids []int) (planets []*Planet) {
	fids = shuffleInts(fids)
	num := 16 * len(fids)
	bigN := num/4 - 1
	littleN := num - (bigN + 1)
	bigArea := bigN * 126
	littleArea := littleN*125 + bigArea
	// a = 3*n*(n+1) + 1
	// n^2 + n = (a-1)/3
	var bigRange, littleRange int
	for i := 1; ; i++ {
		if 3*i*(i+1)+1 > bigArea {
			bigRange = i
			break
		}
	}
	for i := bigRange; ; i++ {
		if 3*i*(i+1)+1 > littleArea {
			littleRange = i
			break
		}
	}
	planets = make([]*Planet, num)
	planets[0] = &Planet{Db: g.Db, Gid: g.Gid, Pid: 999, Name: "Planet Borion", Loc: hexagon.Coord{0, 0}, Inhabitants: 15, Resources: 30}
	names := GetAdj(num)
	usedNums := map[int]bool{0: true}
	usedLocs := map[hexagon.Coord]bool{hexagon.Coord{0, 0}: true}
	for i := 1; i < bigN; i++ {
		p := &Planet{Db: g.Db, Gid: g.Gid, Inhabitants: pick(10), Resources: 10 + pick(10), Name: Nameify(names[i])}
		for usedNums[p.Pid] {
			p.Pid = pick(898) + 99
		}
		for UsedSpace(usedLocs, p.Loc) {
			p.Loc = RandBigPlLoc(bigRange)
		}
		usedNums[p.Pid] = true
		usedLocs[p.Loc] = true
		planets[i] = p
	}
	for i := bigN; i < num-len(fids); i++ {
		p := &Planet{Db: g.Db, Gid: g.Gid, Resources: pick(10), Name: Nameify(names[i])}
		for usedNums[p.Pid] {
			p.Pid = pick(898) + 99
		}
		for UsedSpace(usedLocs, p.Loc) {
			p.Loc = RandLittlePlLoc(bigRange, littleRange)
		}
		usedNums[p.Pid] = true
		usedLocs[p.Loc] = true
		planets[i] = p
	}
	for i := 0; i < len(fids); i++ {
		p := &Planet{Db: g.Db, Gid: g.Gid, Controller: fids[i], Inhabitants: 5, Resources: 15, Parts: 5, Name: Nameify(names[num-1-i])}
		for usedNums[p.Pid] {
			p.Pid = pick(898) + 99
		}
		for UsedSpace(usedLocs, p.Loc) {
			p.Loc = RandHomePlLoc(bigRange, littleRange, i, len(fids))
		}
		usedNums[p.Pid] = true
		usedLocs[p.Loc] = true
		planets[num-1-i] = p
	}
	return planets
}

func UsedSpace(used map[hexagon.Coord]bool, test hexagon.Coord) bool {
	for r := 0; r < 6; r++ {
		for _, cd := range test.Ring(r) {
			if used[cd] {
				return true
			}
		}
	}
	return false
}

func Nameify(str string) string {
	return "Planet " + strings.Title(str)
}

func RandBigPlLoc(bigRange int) hexagon.Coord {
	r := 5 + pick(bigRange-5)
	theta := pick(6*r) - 1
	return hexagon.Polar{r, theta}.Coord()
}

func RandLittlePlLoc(bigRange, littleRange int) hexagon.Coord {
	r := bigRange + pick(littleRange-bigRange)
	theta := pick(6*r) - 1
	return hexagon.Polar{r, theta}.Coord()
}

func RandHomePlLoc(bigRange, littleRange, i, numF int) hexagon.Coord {
	rnge := (littleRange - bigRange) / 2
	start := bigRange + littleRange/4
	r := start + pick(rnge)
	tWidth := 6 * r / numF
	theta := tWidth*i + pick(tWidth) - 1
	return hexagon.Polar{r, theta}.Coord()
}
