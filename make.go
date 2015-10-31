package planetattack

import (
	"database/sql"
)

func MakeGame(db *sql.DB, gameName, owner, password string) (g *Game, err error) {
	g = &Game{db: db, Name: gameName, Owner: owner, Password: password}
	err = g.Insert()
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (g *Game) MakeFaction(uName, fName string) (f *Faction, err error) {
	f = &Faction{db: g.db, Gid: g.Gid, Name: fName, Owner: uName}
	err = f.Insert()
	if err != nil {
		return nil, err
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
	query := PlanetMassInsertQ(planets)
	res, err := g.db.Exec(query)
	if err != nil {
		return Log(err)
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		return Log("failed to insert planets", g.Gid, ": 0 rows affected")
	}
	// VEIWS
	query = ViewMassInsertQ(planets)
	res, err = g.db.Exec(query)
	if err != nil {
		return Log(err)
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		return Log("failed to insert faction views for", g.Gid, ": 0 rows affected")
	}
	//
	query = PlanetViewMassInsertQ(planets, fids)
	res, err = g.db.Exec(query)
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
