package planetattack

import (
	"fmt"
	"mule/hexagon"
	"strings"
)

func (g *Game) RunTurn() {
	orders, err := g.AllOrders()
	if err != nil {
		return
	}
	ships, err := g.AllShips()
	if err != nil {
		return
	}
	planetL, err := g.AllPlanets()
	if err != nil {
		return
	}
	factions, err := g.Factions()
	if err != nil {
		return
	}
	planets := make(map[int]*Planet, len(planetL))
	galaxy := make(map[hexagon.Coord]*Planet, len(planetL))
	viewpoints := map[int][]hexagon.Coord{}
	for _, pl := range planetL {
		planets[pl.Pid] = pl
		galaxy[pl.Loc] = pl
		if pl.Controller == 0 {
			continue
		}
		if list, ok := viewpoints[pl.Controller]; ok {
			viewpoints[pl.Controller] = append(list, pl.Loc)
		} else {
			viewpoints[pl.Controller] = []hexagon.Coord{pl.Loc}
		}
	}
	// FIRE THE ROCKETS!! //
	for _, o := range orders {
		source, target := planets[o.Source], planets[o.Target]
		if source == nil || target == nil || source.Controller != o.Fid || source.Parts < o.Size || source.Pid == target.Pid {
			Log("Bad order:", o, source, target)
			return
		}
		sh := source.LaunchShip(target, o.Fid, o.Size)
		ships = append(ships, sh)
	}
	// TRAVEL THE COSMOS!! //
	landing := map[int][]*Ship{}
	notLanding := []*Ship{}
	var maxDist int
	for _, sh := range ships {
		land, dist := sh.Travel()
		if land {
			if dist > maxDist {
				maxDist = dist
			}
			if list, ok := landing[dist]; ok {
				landing[dist] = append(list, sh)
			} else {
				landing[dist] = []*Ship{sh}
			}
		} else {
			notLanding = append(notLanding, sh)
		}
	}

	// ========  LAND ON PLANETS!!  ========== //
	newViews := map[[2]int]*PlanetView{}
	var landCount int
	for i := 1; i <= maxDist; i++ {
		if len(landing[i]) == 0 {
			continue
		}
		landers := shuffleShips(landing[i])
		for _, sh := range landers {
			landCount++
			target, ok := galaxy[sh.Path[len(sh.Path)-1]]
			if !ok {
				Log("Couldn't find planet for", sh, "to land on!")
				return
			}
			views, reports := sh.LandOn(target, g.Turn)
			for _, v := range views {
				newViews[[2]int{v.Fid, v.Pid}] = v
			}
			_ = reports
		}
	}
	// PLANETS DO THEIR THING //
	for _, pl := range planetL {
		pl.Produce()
		if pl.Controller != 0 {
			newViews[[2]int{pl.Controller, pl.Pid}] = pl.MakeView(g.Turn, pl.Controller)
		}
	}

	// ================== OK BEGIN SQL POINT OF NO RETURN ================= //
	// DELETE ships landing ( if sid != 0 )
	query := "DELETE FROM ships WHERE "
	parts := []string{}
	for _, list := range landing {
		for _, sh := range list {
			if sh.Sid == 0 {
				continue
			}
			parts = append(parts, fmt.Sprintf("(gid = %d AND fid = %d AND sid = %d)", sh.Gid, sh.Fid, sh.Sid))
		}
	}
	if len(parts) > 0 {
		query += strings.Join(parts, " OR ")
		res, err := g.db.Exec(query)
		if err != nil {
			Log(err)
			return
		}
		if aff, err := res.RowsAffected(); err != nil || aff < 1 {
			Log("failed to delete ships", g.Gid, g.Turn, ": 0 rows affected")
			return
		}
	}
	// UPDATE/INSERT ships not landing ( if sid == 0 INSERT )
	query = "UPDATE ships SET loc = $1 WHERE gid = $2 AND fid = $3 AND sid = $4"
	// prep stmt1
	query = "INSERT INTO ships (gid, fid, size, loc, path) VALUES ($1, $2, $3, $4, $5) returning sid"
	for _, sh := range notLanding {
		if sh.Sid == 0 {
			//err := g.db.QueryRow(query, sh.Gid, sh.Fid, sh.Size, sh.Loc, sh.Path).Scan(&(sh.Sid))
		} else {
			// stmt1.Exec
		}
	}
	// DELETE old orders
	query = "DELETE FROM orders WHERE gid = $1"
	// DELETE old shipViews
	query = "DELETE FROM shipviews WHERE gid = $1"
	// INSERT new shipViews
	shipViews := []*ShipView{}
	for _, sh := range ships {
		for fid, _ := range factions {
			shipViews = append(shipViews, sh.MakeView(viewpoints[fid], fid))
		}
	}
	query = "INSERT INTO shipviews (gid, controller, viewer, sid, loc, trail, size) VALUES "
	// UPDATE PLANETVIEWS
	query = "UPDATE planetviews SET turn = $1, controller = $2, inhabitants = $3, resources = $4, parts = $5 WHERE gid = $6 AND fid = $7 AND pid = $8"
	for _, v := range newViews {
		_ = v
	}
	// UPDATE PLANETS
	query = "UPDATE planets SET controller = $1, inhabitants = $2, resources = $3, parts = $4 WHERE gid = $5 AND pid = $6"
	for _, pl := range planets {
		if pl.Modified() {
			//
		}
	}
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
	facs, err := g.Factions()
	if err != nil {
		return err
	}
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
