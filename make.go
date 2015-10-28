package planetattack

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
	// PLANETS //
	planets := g.MakeGalaxy(fids)
	fmt.Println("Made planets:", planets)
	for _, pl := range planets {
		fmt.Printf("%d||", pl.Loc)
	}
	fmt.Println("")
	query := PlanetMassInsertQ(planets)
	fmt.Println("first query:\n", query, "\n")
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
