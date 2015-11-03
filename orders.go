package planetattack

import (
	"database/sql"
	//	"mule/hexagon"
)

type Order struct {
	db     *sql.DB
	Gid    int
	Fid    int
	Source int
	Target int
	Size   int
}

func (f *Faction) OrdersFor(source int) []*Order {
	query := "SELECT target, size FROM orders WHERE gid = $1 AND fid = $2 AND source = $3"
	r := []*Order{}
	rows, err := f.db.Query(query, f.Gid, f.Fid, source)
	if err != nil {
		Log(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		o := &Order{db: f.db, Gid: f.Gid, Fid: f.Fid, Source: source}
		err = rows.Scan(&(o.Target), &(o.Size))
		if err != nil {
			Log("order scan problem: ", err)
			return nil
		}
		r = append(r, o)
	}
	if err = rows.Err(); err != nil {
		Log(err)
		return nil
	}
	return r
}

func (f *Faction) AllOrders() []*Order {
	query := "SELECT source, target, size FROM orders WHERE gid = $1 AND fid = $2"
	r := []*Order{}
	rows, err := f.db.Query(query, f.Gid, f.Fid)
	if err != nil {
		Log(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		o := &Order{db: f.db, Gid: f.Gid, Fid: f.Fid}
		err = rows.Scan(&(o.Source), &(o.Target), &(o.Size))
		if err != nil {
			Log("order scan problem: ", err)
			return nil
		}
		r = append(r, o)
	}
	if err = rows.Err(); err != nil {
		Log(err)
		return nil
	}
	return r
}

func (g *Game) AllOrders() ([]*Order, error) {
	query := "SELECT fid, source, target, size FROM orders WHERE gid = $1"
	r := []*Order{}
	rows, err := g.db.Query(query, g.Gid)
	if err != nil {
		return nil, Log(err)
	}
	defer rows.Close()
	for rows.Next() {
		o := &Order{db: g.db, Gid: g.Gid}
		err = rows.Scan(&(o.Fid), &(o.Source), &(o.Target), &(o.Size))
		if err != nil {
			return nil, Log(err)
		}
		r = append(r, o)
	}
	if err = rows.Err(); err != nil {
		return nil, Log(err)
	}
	return r, nil
}

func (f *Faction) DelOrder(source, target int) error {
	query := "DELETE FROM orders WHERE gid = $1 AND fid = $2 AND source = $3 AND target = $4"
	_, err := f.db.Exec(query, f.Gid, f.Fid, source, target)
	if err != nil {
		Log("failed to delete order", f.Gid, f.Fid, source, target, ":", err)
		return err
	}
	/*
		if aff, err := res.RowsAffected(); err != nil || aff < 1 {
			Log("failed to delete order", f.Gid, f.Fid, source, target, ": 0 rows affected")
			return err
		}*/
	return nil
}

func (g *Game) DelAllOrders() error {
	query := "DELETE FROM orders WHERE gid = $1"
	_, err := g.db.Exec(query, g.Gid)
	if err != nil {
		Log("failed to delete order", g.Gid, ":", err)
		return err
	}
	return nil
}

func (f *Faction) SetOrder(source, target, size int) error {
	if size < 1 {
		return f.DelOrder(source, target)
	}
	query := "UPDATE orders SET size = $1 WHERE gid = $2 AND fid = $3 AND target = $4 AND source = $5;"
	res, err := f.db.Exec(query, size, f.Gid, f.Fid, target, source)
	if err != nil {
		Log("failed to update order", f.Gid, f.Fid, source, target, size, ":", err)
		return err
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		query = "INSERT INTO orders (gid, fid, source, target, size) VALUES ($1, $2, $3, $4, $5)"
		res, err = f.db.Exec(query, f.Gid, f.Fid, source, target, size)
		if err != nil {
			Log("failed to insert order", f.Gid, f.Fid, source, target, size, ":", err)
			return err
		}
		if aff, err = res.RowsAffected(); err != nil || aff < 1 {
			Log("failed to upsert order", f.Gid, f.Fid, source, target, size, ": 0 rows affected")
			return err
		}
	}
	return nil
}
