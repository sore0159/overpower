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

func (f *Faction) SetOrder(source, target, size int) error {
	if size < 1 {
		return f.DelOrder(source, target)
	}
	// ===========  TODO  ================= //
	inQuery := "INSERT INTO orders (gid, fid, source, target, size) VALUES ($1, $2, $3, $4, $5)"
	upQuery := "UPDATE orders SET size = $1 WHERE gid = $2 AND fid = $3 AND source = $4 AND target = $5"
	_, _ = inQuery, upQuery
	return nil
}
