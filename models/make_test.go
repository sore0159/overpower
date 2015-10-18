package attack

import (
	"log"
	"testing"
)

func TestFirst(t *testing.T) {
	log.Println("TESTING")
}

func TestSecond(t *testing.T) {
	db, err := LoadDB()
	if err != nil {
		log.Println("Error getting db")
		return
	}
	log.Println("Got db")
	facs := map[string]string{
		"fac1": "user1",
		"fac2": "user2",
		"fac3": "user3",
	}
	g, err := MakeGame(db, "Game Name", "user1")
	if err != nil {
		log.Println("GAME 1 CONSTRUCTION FAILED")
		return
	}
	log.Println("GAME 1 MADE:", g.Gid)
	var fid int
	var fac *Faction
	for fName, uName := range facs {
		f, err := g.AddFac(fName, uName)
		if err != nil {
			log.Println("FAC", fName, "CONSTRUCTION FAILED")
			return
		}
		log.Println("FAC", f, "MADE")
		fid = f.Fid
		fac = f
	}
	err = g.Start()
	if err != nil {
		log.Println("GAME START FAILED:", err)
		return
	}
	log.Println("GAME 1 STARTED")
	var pl2 *Planet
	for loc, pl := range g.Planets() {
		log.Println("FOUND PLANET AT", loc, ":", pl)
		pl2 = pl
	}
	stmt, err := g.UpdateViewStmt()
	if err != nil {
		log.Println("Error generating update view stmt", err)
	} else {
		err = pl2.UpdateView(stmt, fid, 2)
		if err != nil {
			log.Println("Error updating view", err)
		} else {
			log.Println("Updated View")
		}
	}
	for _, pv := range fac.PlanetViews() {
		log.Println("Found planetview:", pv)
	}
	DelGame(db, g.Gid)
}

func XTestThird(t *testing.T) {
	gid := 38
	db, err := LoadDB()
	if err != nil {
		log.Println("Error getting db")
		return
	}
	log.Println("Got db")
	g := &Game{Db: db, Gid: gid}
	g.Select()
	log.Println("Got", g)
	log.Println("Planets:", g.Planets())
}

func XTestFourth(t *testing.T) {
	str := []byte("(305,5)")
	log.Println(str)
	p := Point{}
	(&p).Scan(str)
	log.Println("POINT:", p)
}
