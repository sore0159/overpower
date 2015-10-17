package models

import (
	"database/sql"
	"fmt"
	"strings"
)

func MakeGame(db *sql.DB, name, owner string, facs map[string]string) Game {
	// GAME //
	var gid int
	query := "INSERT INTO games (name, owner, turn) VALUES($1, $2, $3) RETURNING id"
	err := db.QueryRow(query, name, owner, 1).Scan(&gid)
	if err != nil {
		Log(err)
		return nil
	}
	g := GameDB{db, gid}
	// FACTIONS //
	query = "INSERT INTO factions (gid, owner, name) VALUES ($1, $2, $3) RETURNING fid"
	stm, err := DB.Prepare(query)
	if err != nil {
		Log(err)
		return nil
	}
	var fids []int
	for fName, uName := range facs {
		var fid int
		err = stm.QueryRow(gid, uName, fName).Scan(&fid)
		if err != nil {
			Log(err)
			stm.Close()
			return nil
		}
		fids = append(fids, fid)
	}
	stm.Close()
	// PLANETS //
	query = "INSERT INTO planets (gid, pid, name, loc) VALUES "
	queryParts := []string{}
	viewQuery := "INSERT INTO planetviews (gid, fid, pid, name, loc, turn) VALUES "
	vQueryParts := []string{}
	fNum := len(facs)
	pNum := fNum * 20
	names := shuffleWords(GetAdj())
	usedNums := map[int]bool{0: true}
	usedLocs := map[[2]int]bool{}
	for i := 0; i < pNum; i++ {
		name = names[i]
		var pid int
		for usedNums[pid] {
			pid = pick(898) + 99
		}
		usedNums[pid] = true
		var loc [2]int
		for UsedSpace(usedLocs, loc) {
			loc = [2]int{pick(200) - 100, pick(200) - 100}
		}
		usedLocs[loc] = true
		queryParts = append(queryParts, fmt.Sprintf("(%d, %d, '%s', POINT(%d, %d))", gid, pid, name, loc[0], loc[1]))
		for _, fid := range fids {
			vQueryParts = append(vQueryParts, fmt.Sprintf("(%d, %d, %d, '%s', POINT(%d, %d), 0)", gid, fid, pid, name, loc[0], loc[1]))
		}
	}
	query += strings.Join(queryParts, ", ")
	stm, err = DB.Prepare(query)
	if err != nil {
		Log(query)
		Log(err)
		return nil
	}
	defer stm.Close()
	viewQuery += strings.Join(vQueryParts, ", ")
	stmV, err := DB.Prepare(viewQuery)
	if err != nil {
		Log(viewQuery)
		Log(err)
		return nil
	}
	defer stmV.Close()
	res, err := stm.Exec()
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
	res, err = stmV.Exec()
	if err != nil {
		Log(err)
		return nil
	}
	if aff, err := res.RowsAffected(); err != nil {
		Log(err)
		return nil
	} else if aff == 0 {
		Log("Planet view creation failed")
		return nil
	}

	return g
}

func UsedSpace(used map[[2]int]bool, test [2]int) bool {
	_, ok := used[test]
	return ok
}

func DelGame(db *sql.DB, id int) {
	query := "DELETE FROM games where id = $1"
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
