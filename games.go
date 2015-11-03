package planetattack

import (
	"database/sql"
	"errors"
	"fmt"
	"mule/hexagon"
	"strings"
)

type Game struct {
	db       *sql.DB
	Gid      int
	Owner    string
	Name     string
	Turn     int
	Password string
	//
	CacheFactions map[int]*Faction
	CachePlanets  map[hexagon.Coord]*Planet
	CacheShips    []*Ship
}

func NewGame(db *sql.DB) *Game {
	return &Game{
		db: db,
	}
}

func GetGames(db *sql.DB, gids []int) []*Game {
	if len(gids) < 1 {
		return []*Game{}
	}
	query := "SELECT gid, owner, name, turn FROM games WHERE "
	parts := make([]string, len(gids))
	games := make([]*Game, len(gids))
	for i, gid := range gids {
		parts[i] = fmt.Sprintf("gid = %d", gid)
	}
	query += strings.Join(parts, " OR ")
	rows, err := db.Query(query)
	if err != nil {
		Log(err)
		return nil
	}
	defer rows.Close()
	var i int
	for rows.Next() {
		g := &Game{db: db}
		err = rows.Scan(&(g.Gid), &(g.Owner), &(g.Name), &(g.Turn))
		if err != nil {
			Log("game scan problem: ", err)
			return nil
		}
		if i > len(games)-1 {
			Log("game scan problem: too many entries!", i, len(games))
			return nil
		}
		games[i] = g
		i++
	}
	if err = rows.Err(); err != nil {
		Log(err)
		return nil
	}
	return games
}

func (g *Game) Insert() error {
	query := "INSERT INTO games (name, owner, turn, password) VALUES($1, $2, $3, $4) RETURNING gid"
	var pswd sql.NullString
	if g.Password != "" {
		pswd.Valid = true
		pswd.String = g.Password
	}
	err := g.db.QueryRow(query, g.Name, g.Owner, g.Turn, pswd).Scan(&(g.Gid))
	if err != nil {
		return Log(err)
	}
	return nil
}

func (g *Game) Select() bool {
	var err error
	var pswd sql.NullString
	if g.Gid != 0 {
		query := "SELECT owner, name, turn, password FROM games WHERE gid = $1"
		err = g.db.QueryRow(query, g.Gid).Scan(&(g.Owner), &(g.Name), &(g.Turn), &pswd)
	} else if g.Owner != "" {
		query := "SELECT gid, name, turn, password FROM games WHERE owner = $1"
		err = g.db.QueryRow(query, g.Owner).Scan(&(g.Gid), &(g.Name), &(g.Turn), &pswd)
	} else {
		err = errors.New("tried to SELECT game with no gid/owner")
	}
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		Log("GAME SELECT ERROR FOR GAME", g, ":", err)
		return false
	}
	if pswd.Valid {
		g.Password = pswd.String
	}
	return true
}

func (g *Game) IncTurn() error {
	g.Turn += 1
	query := "UPDATE games SET turn = turn + 1 WHERE gid = $1"
	res, err := g.db.Exec(query, g.Gid)
	if err != nil {
		return Log("failed to inc game", g.Gid, "turn:", err)
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		return Log("failed to inc game", g.Gid, "turn: 0 rows affected")
	}
	return nil
}

func (g *Game) Factions() (map[int]*Faction, error) {
	if g.CacheFactions == nil {
		g.CacheFactions = map[int]*Faction{}
		query := "SELECT fid, owner, name, done FROM factions WHERE gid = $1"
		rows, err := g.db.Query(query, g.Gid)
		if err != nil {
			g.CacheFactions = nil
			return nil, Log(err)
		}
		defer rows.Close()
		for rows.Next() {
			f := &Faction{db: g.db, Gid: g.Gid}
			err = rows.Scan(&(f.Fid), &(f.Owner), &(f.Name), &(f.Done))
			if err != nil {
				g.CacheFactions = nil
				return nil, Log(err)
			}
			g.CacheFactions[f.Fid] = f
		}
		if err = rows.Err(); err != nil {
			g.CacheFactions = nil
			return nil, Log(err)
		}
	}
	return g.CacheFactions, nil
}

func (g *Game) AllPlanets() ([]*Planet, error) {
	r := []*Planet{}
	query := "SELECT pid, name, loc, controller, inhabitants, resources, parts FROM planets WHERE gid = $1"
	rows, err := g.db.Query(query, g.Gid)
	if err != nil {
		return nil, Log(err)
	}
	defer rows.Close()
	for rows.Next() {
		p := &Planet{db: g.db, Gid: g.Gid}
		var controller sql.NullInt64
		err = rows.Scan(&(p.Pid), &(p.Name), &(p.Loc), &controller, &(p.Inhabitants), &(p.Resources), &(p.Parts))
		if err != nil {
			return nil, Log(err)
		}
		if controller.Valid {
			p.Controller = int(controller.Int64)
		}
		r = append(r, p)
	}
	if err = rows.Err(); err != nil {
		return nil, Log(err)
	}
	return r, nil
}

func (g *Game) GetPlanets(pids ...int) []*Planet {
	if len(pids) < 1 {
		return nil
	}
	r := make([]*Planet, len(pids))
	query := "SELECT pid, name, loc, controller, inhabitants, resources, parts FROM planets WHERE gid = $1 AND ("
	parts := make([]string, len(pids))
	mp := make(map[int]int, len(pids))
	for i, pid := range pids {
		mp[pid] = i
		parts[i] = fmt.Sprintf("pid = %d", pid)
	}
	query += strings.Join(parts, " OR ") + ")"
	rows, err := g.db.Query(query, g.Gid)
	if err != nil {
		Log(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		p := &Planet{db: g.db, Gid: g.Gid}
		var controller sql.NullInt64
		err = rows.Scan(&(p.Pid), &(p.Name), &(p.Loc), &controller, &(p.Inhabitants), &(p.Resources), &(p.Parts))
		if err != nil {
			Log(err)
			return nil
		}
		if controller.Valid {
			p.Controller = int(controller.Int64)
		}
		r[mp[p.Pid]] = p
	}
	if err = rows.Err(); err != nil {
		Log(err)
		return nil
	}
	return r
}

func (g *Game) AllShips() ([]*Ship, error) {
	if g.CacheShips == nil {
		g.CacheShips = []*Ship{}
		query := "SELECT fid, sid, size, loc, path WHERE gid = $1"
		rows, err := g.db.Query(query, g.Gid)
		if err != nil {
			g.CacheShips = nil
			return nil, Log(err)
		}
		defer rows.Close()
		for rows.Next() {
			s := &Ship{db: g.db, Gid: g.Gid}
			err = rows.Scan(&(s.Fid), &(s.Sid), &(s.Size), &(s.Loc), &(s.Path))
			if err != nil {
				g.CacheShips = nil
				return nil, Log(err)
			}
			g.CacheShips = append(g.CacheShips, s)
		}
		if err = rows.Err(); err != nil {
			g.CacheShips = nil
			return nil, Log(err)
		}
	}
	return g.CacheShips, nil
}
