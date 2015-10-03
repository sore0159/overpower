package models

//import "fmt"

func MakeGame(name, owner string, facs map[string]string) *Game {
	g := &Game{}
	g.Name = name
	g.Owner = owner
	g.Turn = 1
	query := "INSERT INTO games (name, owner, turn) VALUES($1, $2, $3) RETURNING id"
	err := DB.QueryRow(query, name, owner, g.Turn).Scan(&(g.GID))
	if err != nil {
		Log(err)
		return nil
	}
	g.Factions = make(map[int]*Faction, len(facs))
	query = "INSERT INTO factions (gid, owner, name) VALUES ($1, $2, $3) RETURNING fid"
	stm, err := DB.Prepare(query)
	if err != nil {
		Log(err)
		return nil
	}
	for fName, uName := range facs {
		f := NewFaction()
		f.GID = g.GID
		f.Owner = uName
		f.Name = fName
		err = stm.QueryRow(g.GID, uName, fName).Scan(&(f.FID))
		if err != nil {
			Log(err)
			stm.Close()
			return nil
		}
		g.Factions[f.FID] = f
	}
	stm.Close()
	query = "INSERT INTO planets (pid, gid, name, loc) VALUES ($1, $2, $3, POINT($4, $5) )"
	stm, err = DB.Prepare(query)
	if err != nil {
		Log(err)
		return nil
	}
	defer stm.Close()
	g.MakeSector()
	for _, pl := range g.Sector.Planets {
		//fmt.Println("Making planet", pl)
		res, err := stm.Exec(pl.PID, g.GID, pl.Name, pl.Loc[0], pl.Loc[1])
		if err != nil {
			Log(err)
			return nil
		}
		if aff, err := res.RowsAffected(); err != nil {
			Log(err)
			return nil
		} else if aff == 0 {
			Log("Planet creation failed")
			return nil
		}
	}
	return g
}

func DelGame(id int) {
	query := "DELETE FROM games where id = $1"
	res, err := DB.Exec(query, id)
	if err != nil {
		Log("failed to delete game", id, ":", err)
		return
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		Log("failed to delete game", id, ": 0 rows affected")
		return
	}
}

func (g *Game) MakeSector() {
	fNum := len(g.Factions)
	g.Sector = NewSector()
	pNum := fNum * 20
	names := shuffleWords(GetAdj())
	usedNums := map[int]bool{0: true}
	for i := 0; i < pNum; i++ {
		pl := NewPlanet()
		pl.Name = names[i]
		var pID int
		for usedNums[pID] {
			pID = pick(898) + 99
		}
		pl.GID = g.GID
		pl.PID = pID
		usedNums[pl.PID] = true
		var loc [2]int
		for !g.Sector.Clear(loc) {
			loc = [2]int{pick(200) - 100, pick(200) - 100}
		}
		pl.Loc = loc
		g.Sector.Planets[loc] = pl
	}
}

func (s *Sector) Clear(loc [2]int) bool {
	_, ok := s.Planets[loc]
	return !ok
}
