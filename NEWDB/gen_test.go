package db

import (
	"log"
	"testing"
)

func Test1(t *testing.T) {
	log.Println("TEST ONE")
}

func Test2(t *testing.T) {
	db, err := LoadDB()
	if my, bad := Check(err, "test2 loaddb fail"); bad {
		log.Println(my.MuleError())
		return
	}
	log.Println("DB loaded")
	games, err := db.GetGames(C{"gid", 1})
	if my, bad := Check(err, "test2 fail", "games", games); bad {
		log.Println(my.MuleError())
		return
	}
	log.Println("Got games", len(games))
	game, err := db.GetGame(C{"owner", "AutoTest"})
	if my, bad := Check(err, "test2 getgame fail", "owner", "AutoTest"); bad {
		if !my.BaseIs(ErrNoneFound) {
			log.Println(my.MuleError())
			return
		}
		err = db.MakeGame("AutoTest", "Automade", "", 40)
		if my, bad := Check(err, "test2 automake game fail"); bad {
			log.Println(my.MuleError())
			return
		}
		log.Println("Made auto-test game")
		game, err = db.GetGame(C{"owner", "AutoTest"})
		if my, bad := Check(err, "test2 getgame fail two", "owner", "AutoTest"); bad {
			log.Println(my.MuleError())
			return
		}
	}
	log.Println("Looking at game:", game)
	gid := game.Gid()
	// -------------------- //
	game.IncTurn()
	log.Println("Incing Turn to:", game.Turn())
	err = db.UpdateGames(game)
	if my, bad := Check(err, "test2 update failure"); bad {
		log.Println(my.MuleError())
		return
	}
	log.Println("Updated game", gid)

	err = db.DropGames(C{"gid", gid})
	if my, bad := Check(err, "test2 failed dropgame", "gid", gid); bad {
		log.Println(my.MuleError())
		return
	}
	log.Println("Dropped game", gid)
}
