package planetattack

import (
	"database/sql"
)

type Request struct {
	db    *sql.DB
	Gid   int
	Owner string
	Name  string
}

func (r Request) Insert() error {
	query := "INSERT INTO requests (gid, owner, name) VALUES ($1, $2, $3)"
	res, err := r.db.Exec(query, r.Gid, r.Owner, r.Name)
	if err != nil {
		return Log(err)
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		return Log("failed to insert request", r, ": 0 rows affected")
	}
	return nil
}

func (r Request) Approve() (f *Faction, err error) {
	f = &Faction{db: r.db, Gid: r.Gid, Name: r.Name, Owner: r.Owner}
	err = f.Insert()
	if err != nil {
		return nil, err
	}
	r.Delete()
	return f, nil
}

func (r Request) Delete() {
	query := "DELETE FROM requests where gid = $1 AND owner = $2"
	res, err := r.db.Exec(query, r.Gid, r.Owner)
	if err != nil {
		Log("failed to delete request", r.Gid, r.Owner, ":", err)
		return
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		Log("failed to delete request", r.Gid, r.Owner, ": 0 rows affected")
		return
	}
}

func (g *Game) AllRequests() []Request {
	query := "SELECT owner, name FROM requests WHERE gid = $1"
	r := make([]Request, 0)
	rows, err := g.db.Query(query, g.Gid)
	if err != nil {
		Log(err)
		return r
	}
	defer rows.Close()
	for rows.Next() {
		rq := Request{db: g.db, Gid: g.Gid}
		err = rows.Scan(&(rq.Owner), &(rq.Name))
		if err != nil {
			Log(err)
			return []Request{}
		}
		r = append(r, rq)
	}
	if err = rows.Err(); err != nil {
		Log(err)
		return []Request{}
	}
	return r
}

func AllRequests(db *sql.DB, userN string) []Request {
	query := "SELECT gid, name FROM requests WHERE owner = $1"
	r := make([]Request, 0)
	rows, err := db.Query(query, userN)
	if err != nil {
		Log(err)
		return r
	}
	defer rows.Close()
	for rows.Next() {
		rq := Request{db: db, Owner: userN}
		err = rows.Scan(&(rq.Gid), &(rq.Name))
		if err != nil {
			Log(err)
			return []Request{}
		}
		r = append(r, rq)
	}
	if err = rows.Err(); err != nil {
		Log(err)
		return []Request{}
	}
	return r
}

func (g *Game) GetRequest(userN string) (r Request, ok bool) {
	query := "SELECT name FROM requests WHERE gid = $1 AND owner = $2"
	r.db = g.db
	r.Owner = userN
	r.Gid = g.Gid
	err := r.db.QueryRow(query, r.Gid, r.Owner).Scan(&(r.Name))
	if err != nil {
		return r, false
	}
	return r, true
}
