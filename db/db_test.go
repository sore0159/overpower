package db

import (
	"fmt"
	"mule/hexagon"
	"mule/mydb"
	"testing"
)

func TestFirst(t *testing.T) {
	fmt.Println("TESTING")
}

func TestSecond(t *testing.T) {
	db, ok := LoadDB()
	fmt.Println("TEST SECOND Got DB ok:", ok)
	if !ok {
		return
	}
	g, ok := db.GetGame(1)
	if !ok {
		fmt.Println("FAILED GETGAME")
		game := &Game{gid: 1}
		if !game.Insert(db.db) {
			fmt.Println("FAILED INSERT")
			return
		}
		fmt.Println("INSERTED")
		g, ok = db.GetGame(1)
		if !ok {
			fmt.Println("FAILED SECOND GETGAME")
			return
		}
	}
	fmt.Println("GOTGAME")
	fmt.Println("TURN:", g.Turn())
	g.IncTurn()
	updateList := []mydb.Updater{}
	updateList = append(updateList, g)
	ok = db.Update(updateList)
	if !ok {
		fmt.Println("UPDATE FAILED!")
		return
	}
	fmt.Println("UPDATE PASSED")
	f, ok := db.GetFidFaction(g.Gid(), 1)
	if !ok {
		fmt.Println("GET FAC 1 FAILED")
		ok = db.MakeFaction(g.Gid(), "TEST1", "TESTERS")
		if !ok {
			fmt.Println("MAKE FACTION1 FAILED")
			return
		}
		f, ok = db.GetFidFaction(g.Gid(), 1)
		if !ok {
			fmt.Println("FAILED GET2 FAC1")
			return
		}

	}
	fmt.Println("GOT FACTION", f)
	ships, ok := db.GetAllGidShips(g.Gid())
	if !ok {
		fmt.Println("FAILED GET SHIPS")
		return
	}
	for _, s := range ships {
		fmt.Println("GOT SHIP:", s)
	}
	if len(ships) > 0 {
		s, _ := ships[0].(*Ship)
		s.path = []hexagon.Coord{hexagon.Coord{0, 1}, {0, 2}}
		ok = s.Insert(db.db)
	} else {
		fmt.Println("NO SHIPS FOUND")
		s := &Ship{1, 1, 0, 5, 5, []hexagon.Coord{{0, 1}}}
		ok = s.Insert(db.db)
	}
	if !ok {
		fmt.Println("INSERT FAILED")
		return
	} else {
		fmt.Println("INSERT PASSED")
	}
	ships, ok = db.GetAllGidShips(g.Gid())
	if !ok {
		fmt.Println("FAILED GET SHIPS 2")
		return
	}
	for _, s := range ships {
		fmt.Println("GOT2 SHIP:", s)
	}
	shipViews, ok := db.GetFidTurnShipViews(g.Gid(), 1, 1)
	if !ok {
		fmt.Println("FAILED GETSHIPVIEWS")
		return
	}
	if len(shipViews) < 1 {
		fmt.Println("GOT NO SHIPVIEWS")
		at := hexagon.Coord{1, 1}
		seen := []hexagon.Coord{{0, 0}, {1, 0}, {1, 1}}
		for _, s := range ships {
			_, _ = at, seen
			fmt.Println("TODO: shipview from", s)
			/*
				sv := MakeShipView(s, 1, 1, at, seen)
				if !sv.Insert(db.db) {
					fmt.Println("FAILED SHIPVIEW INSERT")
					return
				}
			*/
		}
		shipViews, ok = db.GetFidTurnShipViews(g.Gid(), 1, 1)
		if len(shipViews) < 1 {
			fmt.Println("GOT2 NO SHIPVIEWS")
			return
		}
	}
	for _, sv := range shipViews {
		fmt.Println("GOT SHIPVIEW", sv)
	}
}
