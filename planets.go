package planetattack

import (
	"database/sql"
	"fmt"
	"mule/hexagon"
	"strings"
)

type Planet struct {
	db          *sql.DB
	Gid         int
	Pid         int
	Name        string
	Loc         hexagon.Coord
	Controller  int
	Inhabitants int
	Resources   int
	Parts       int
	conMod      bool
	inhMod      bool
	resMod      bool
	partsMod    bool
	//
	Arrivals int
}

func (p *Planet) Modified() bool {
	return p.conMod || p.inhMod || p.resMod || p.partsMod
}

func (p *Planet) Produce() {
	if p.Controller != 0 && p.Inhabitants > 0 && p.Resources > 0 {
		amount := p.Inhabitants
		if amount > p.Resources {
			amount = p.Resources
		}
		p.SetResources(p.Resources - amount)
		p.SetParts(p.Parts + amount)
	}
	if p.Arrivals > 0 {
		p.SetInhabitants(p.Inhabitants + p.Arrivals)
	}
}

func (p *Planet) MakeView(turn, fid int) *PlanetView {
	return &PlanetView{
		db:          p.db,
		Gid:         p.Gid,
		Fid:         fid,
		Pid:         p.Pid,
		Name:        p.Name,
		Loc:         p.Loc,
		Turn:        turn,
		Controller:  p.Controller,
		Inhabitants: p.Inhabitants + p.Arrivals,
		Resources:   p.Resources,
		Parts:       p.Parts,
	}
}

func (p *Planet) SetController(fid int) {
	if fid != p.Controller {
		p.Controller = fid
		p.conMod = true
	}
}

func (p *Planet) SetInhabitants(num int) {
	if p.Inhabitants != num {
		p.Inhabitants = num
		p.inhMod = true
	}
}

func (p *Planet) SetResources(num int) {
	if p.Resources != num {
		p.Resources = num
		p.resMod = true
	}
}

func (p *Planet) SetParts(num int) {
	if p.Parts != num {
		p.Parts = num
		p.partsMod = true
	}
}

func (p *Planet) Select() error {
	query := "SELECT name, loc, controller, inhabitants, resources, parts FROM planets WHERE gid = $1 AND pid = $2"
	var controller sql.NullInt64
	err := p.db.QueryRow(query, p.Gid, p.Pid).Scan(&(p.Name), &(p.Loc), &controller, &(p.Inhabitants), &(p.Resources), &(p.Parts))
	if err != nil {
		return Log(err)
	}
	if controller.Valid {
		p.Controller = int(controller.Int64)
	}
	return nil
}

func (*Planet) InsertQStart() string {
	return "INSERT INTO planets VALUES (gid, pid, name, loc, controller, inhabitants, resources, parts) "

}

func (p *Planet) InsertQVals() string {
	if p.Controller == 0 {
		return fmt.Sprintf("(%d, %d, '%s', POINT(%d,%d), NULL, %d, %d, %d)", p.Gid, p.Pid, p.Name, p.Loc[0], p.Loc[1], p.Inhabitants, p.Resources, p.Parts)
	} else {
		return fmt.Sprintf("(%d, %d, '%s', POINT(%d,%d), %d, %d, %d, %d)", p.Gid, p.Pid, p.Name, p.Loc[0], p.Loc[1], p.Controller, p.Inhabitants, p.Resources, p.Parts)
	}

}

func (p *Planet) ViewInsertQVals(fid int) string {
	if fid == p.Controller {
		return fmt.Sprintf("(%d, %d, %d, '%s', POINT(%d,%d), 1, %d, %d, %d, %d)", p.Gid, fid, p.Pid, p.Name, p.Loc[0], p.Loc[1], p.Controller, p.Inhabitants, p.Resources, p.Parts)
	} else {
		return fmt.Sprintf("(%d, %d, %d, '%s', POINT(%d,%d), 0, NULL, NULL, NULL, NULL)", p.Gid, fid, p.Pid, p.Name, p.Loc[0], p.Loc[1])
	}
}

func (g *Game) UpdateViewStmt() (*sql.Stmt, error) {
	query := "UPDATE planetviews SET turn = $1, controller = $2, inhabitants = $3, resources = $4, parts = $5 WHERE gid = $6 AND fid = $7 AND pid = $8"
	stmt, err := g.db.Prepare(query)
	if err != nil {
		return nil, Log(err)
	}
	return stmt, nil
}

func (p *Planet) UpdateView(stmt *sql.Stmt, fid, turn int) error {
	controller := sql.NullInt64{}
	if p.Controller == 0 {
		controller.Valid = false
	} else {
		controller.Valid = true
		controller.Int64 = int64(p.Controller)
	}
	res, err := stmt.Exec(turn, controller, p.Inhabitants, p.Resources, p.Parts, p.Gid, fid, p.Pid)
	if err != nil {
		return Log("failed to update view", p.Gid, p.Pid, fid, ":", err)
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		return Log("failed to update view", p.Gid, p.Pid, fid, ": No rows affected")
	}
	return nil
}

func PlanetMassInsertQ(pls []*Planet) string {
	q := "INSERT INTO planets (gid, pid, name, loc, controller, inhabitants, resources, parts) VALUES "
	qStr := make([]string, len(pls))
	for i, pl := range pls {
		qStr[i] = pl.InsertQVals()
	}
	q += strings.Join(qStr, ", ")
	return q
}

func PlanetViewMassInsertQ(pls []*Planet, fids []int) string {
	q := "INSERT INTO planetviews (gid, fid, pid, name, loc, turn, controller, inhabitants, resources, parts) VALUES "
	allStr := make([]string, len(fids))
	for j, fid := range fids {
		qStr := make([]string, len(pls))

		for i, pl := range pls {
			qStr[i] = pl.ViewInsertQVals(fid)
		}
		allStr[j] = strings.Join(qStr, ", ")
	}
	q += strings.Join(allStr, ", ")
	return q
}
