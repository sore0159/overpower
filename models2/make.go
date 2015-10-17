package attack

import (
	"database/sql"
	"fmt"
	"strings"
)

func MakeGame(db *sql.DB, gameName, owner string, facs map[string]string) (gid int, err error) {
	num := 16 * len(facs)
	bigN := num/4 - 1
	littleN := num - (bigN + 1)
	bigArea := bigN * 126
	littleArea := littleN*125 + bigArea
	// a = 3*n*(n+1) + 1
	// n^2 + n = (a-1)/3
	var bigRange, littleRange int
	for i := 1; ; i++ {
		if 3*i*(i+1)+1 > bigArea {
			bigRange = i
			break
		}
	}
	for i := bigRange; ; i++ {
		if 3*i*(i+1)+1 > littleArea {
			littleRange = i
			break
		}
	}
	// GAME //
	query := "INSERT INTO games (name, owner, size, turn) VALUES($1, $2, $3, $4) RETURNING gid"
	err = db.QueryRow(query, gameName, owner, littleRange, 1).Scan(&gid)
	if err != nil {
		Log(err)
		return
	}
	// FACTIONS //
	query = "INSERT INTO factions (gid, owner, name, done) VALUES ($1, $2, $3, false) RETURNING fid"
	stm, err := db.Prepare(query)
	if err != nil {
		Log(err)
		return
	}
	var fids []int
	for fName, uName := range facs {
		var fid int
		err = stm.QueryRow(gid, uName, fName).Scan(&fid)
		if err != nil {
			Log(err)
			stm.Close()
			return
		}
		fids = append(fids, fid)
	}
	stm.Close()
	// PLANETS //
	query = "INSERT INTO planets (gid, pid, name, loc, controller, inhabitants, resources, parts) VALUES "
	queryParts := []string{}
	viewQuery := "INSERT INTO planetviews (gid, fid, pid, name, loc, turn, controller, inhabitants, resources, parts) VALUES "
	vQueryParts := []string{}
	fNum := len(facs)
	pNum := fNum * 16
	names := GetAdj(pNum)
	usedNums := map[int]bool{0: true}
	usedLocs := map[[2]int]bool{}
	var name string
	var pid, inhab, resources, parts int
	var loc [2]int
	for i := 0; i < pNum; i++ {
		if i == 0 {
			name = "Planet Borion"
			pid = 999
			inhab = 15
			resources = 30
			loc = [2]int{0, 0}
		} else {
			name = "Planet " + names[i]
			for usedNums[pid] {
				pid = pick(898) + 99
			}
			if i < bigN+1 {
				for UsedSpace(usedLocs, loc) {
					loc = [2]int{pick(200) - 100, pick(200) - 100}
				}
				inhab = pick(10)
				resources = 10 + pick(10)
			} else if i < pNum-len(facs) {
				resources = pick(10)
			} else {
				inhab = 5
				parts = 5
				resources = 15
			}
		}
		usedNums[pid] = true
		usedLocs[loc] = true
		if i < pNum-len(facs) {
			queryParts = append(queryParts, fmt.Sprintf("(%d, %d, '%s', POINT(%d, %d), NULL, %d, %d, %d)", gid, pid, name, loc[0], loc[1], inhab, resources, parts))
			for _, fid := range fids {
				vQueryParts = append(vQueryParts, fmt.Sprintf("(%d, %d, %d, '%s', POINT(%d, %d), 0, NULL, NULL, NULL, NULL)", gid, fid, pid, name, loc[0], loc[1]))
			}
		} else {
			ownedID := fids[pNum-i]
			queryParts = append(queryParts, fmt.Sprintf("(%d, %d, '%s', POINT(%d, %d), %d, %d, %d, %d)", gid, pid, name, loc[0], loc[1], ownedID, inhab, resources, parts))
			for _, fid := range fids {
				if fid == ownedID {
					vQueryParts = append(vQueryParts, fmt.Sprintf("(%d, %d, %d, '%s', POINT(%d, %d), 1, %d, %d, %d, %d)", gid, fid, pid, name, loc[0], loc[1], ownedID, inhab, resources, parts))
				} else {
					vQueryParts = append(vQueryParts, fmt.Sprintf("(%d, %d, %d, '%s', POINT(%d, %d), 0, NULL, NULL, NULL, NULL)", gid, fid, pid, name, loc[0], loc[1]))
				}
			}

		}
	}
	query += strings.Join(queryParts, ", ")
	stm, err = db.Prepare(query)
	if err != nil {
		Log(query)
		Log(err)
		return
	}
	defer stm.Close()
	viewQuery += strings.Join(vQueryParts, ", ")
	stmV, err := db.Prepare(viewQuery)
	if err != nil {
		Log(viewQuery)
		Log(err)
		return
	}
	defer stmV.Close()
	res, err := stm.Exec()
	if err != nil {
		Log(err)
		return
	}
	if aff, err := res.RowsAffected(); err != nil {
		Log(err)
		return 0, err
	} else if aff == 0 {
		return 0, Log("Planet creation failed")
	}
	res, err = stmV.Exec()
	if err != nil {
		Log(err)
		return
	}
	if aff, err := res.RowsAffected(); err != nil {
		Log(err)
		return 0, err
	} else if aff == 0 {
		return 0, Log("Planet view creation failed")
	}
	return gid, nil
}

func UsedSpace(used map[[2]int]bool, test [2]int) bool {
	_, ok := used[test]
	return ok
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
