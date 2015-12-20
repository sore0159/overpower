package db

import (
	"fmt"
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
}
