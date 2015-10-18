package attack

import (
	"database/sql"
	"fmt"
)

func MakeGame(db *sql.DB, gameName, owner string) (g *Game, err error) {
	g = &Game{Db: db, Name: gameName, Owner: owner}
	err = g.Insert()
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (g *Game) AddFac(fName, uName string) (f *Faction, err error) {
	f = &Faction{Db: g.Db, Gid: g.Gid, Name: fName, Owner: uName}
	err = f.Insert()
	if err != nil {
		return nil, err
	}
	if g.CacheFactions != nil {
		g.CacheFactions[f.Fid] = f
	}
	return f, nil
}

func (g *Game) Start() (err error) {
	facs := g.Factions()
	fids := make([]int, len(facs))
	var i int
	for fid, _ := range facs {
		fids[i] = fid
		i++
	}
	num := 16 * len(facs)
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
	// PLANETS //
	planets := make([]*Planet, num)
	planets[0] = &Planet{Db: g.Db, Gid: g.Gid, Pid: 999, Name: "Planet Borion", Loc: [2]int{0, 0}, Inhabitants: 15, Resources: 30}
	names := GetAdj(num)
	usedNums := map[int]bool{0: true}
	usedLocs := map[[2]int]bool{[2]int{0, 0}: true}
	for i := 1; i < bigN; i++ {
		p := &Planet{Db: g.Db, Gid: g.Gid, Inhabitants: pick(10), Resources: 10 + pick(10), Name: "Planet " + names[i]}
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
	for i := bigN; i < num-len(facs); i++ {
		p := &Planet{Db: g.Db, Gid: g.Gid, Resources: pick(10), Name: "Planet " + names[i]}
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
	for i := 0; i < len(facs); i++ {
		p := &Planet{Db: g.Db, Gid: g.Gid, Controller: fids[i], Inhabitants: 5, Resources: 15, Parts: 5, Name: "Planet " + names[num-1-i]}
		for usedNums[p.Pid] {
			p.Pid = pick(898) + 99
		}
		for UsedSpace(usedLocs, p.Loc) {
			p.Loc = RandHomePlLoc(bigRange, littleRange, i, len(facs))
		}
		usedNums[p.Pid] = true
		usedLocs[p.Loc] = true
		planets[num-1-i] = p
	}
	fmt.Println("Made planets:", planets, num, littleRange, bigRange, len(facs))
	for _, pl := range planets {
		fmt.Printf("%d||", pl.Loc)
	}
	fmt.Println("")
	query := PlanetMassInsertQ(planets)
	res, err := g.Db.Exec(query)
	if err != nil {
		return Log(err)
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		return Log("failed to insert planets", g.Gid, ": 0 rows affected")
	} else {
		fmt.Println(aff, "rows affected in planets make")
	}
	query = PlanetViewMassInsertQ(planets, fids)
	res, err = g.Db.Exec(query)
	if err != nil {
		return Log(err)
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		return Log("failed to insert planetsviews", g.Gid, ": 0 rows affected")
	}
	// GAME //
	err = g.IncTurn()
	if err != nil {
		return err
	}
	return nil
}

func UsedSpace(used map[[2]int]bool, test [2]int) bool {
	_, ok := used[test]
	return ok
}

func DelGame(db *sql.DB, id int) {
	query := "DELETE FROM games where gid = $1"
	res, err := db.Exec(query, id)
	if err != nil {
		Log("failed to delete game", id, ":", err)
		return
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		Log("failed to delete game", id, ": 0 rows affected")
		return
	}
}
