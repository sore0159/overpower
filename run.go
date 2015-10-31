package planetattack

func (g *Game) RunTurn() {
	// Select orders
	// Clear orders, shipviews
	// Select planets
	// Select ships
	// Build new ships
	// Move ships
	// Create shipviews
	// Land ships
	// Insert ships, shipviews
	// Create reports
	// Planet part creation
	// Update Planets, Planetviews
	g.SetFacsNotDone()
	g.IncTurn()
}

func (g *Game) SetFacsNotDone() error {
	query := "UPDATE factions SET done = FALSE WHERE gid = $1"
	res, err := g.db.Exec(query, g.Gid)
	if err != nil {
		return Log("failed to set all facs done", g.Gid, ":", err)
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		return Log("failed to set all facs done", g.Gid, ": no rows affected")
	}
	if g.CacheFactions != nil {
		for _, f := range g.CacheFactions {
			f.Done = false
		}
	}
	return nil
}

func (g *Game) SetDone(fid int, done bool) error {
	facs := g.Factions()
	if len(facs) == 0 {
		return Log("can't toggledone: no factions found")
	}
	var gorun = true
	for _, f := range facs {
		if f.Fid == fid {
			err := f.SetDone(done)
			if err != nil {
				return err
			}

		}
		if !f.Done {
			gorun = false
		}
	}
	if gorun {
		g.RunTurn()
	}
	return nil
}
