package planetattack

import (
	"log"
	"mule/hexagon"
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
		rq := Request{g.db, g.Gid, uName, fName}
		rq.Insert()
		f, err := rq.Approve()
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
	query := "SELECT name FROM planets WHERE loc ~= $1"
	loc := hexagon.Coord{0, 0}
	var name string
	err = db.QueryRow(query, loc).Scan(&name)
	if err != nil {
		log.Println("Error fetching planet at", loc, ":", err)
	} else {
		log.Println("Found planet:", name)
	}
	factions := AllFactions(db, "user2")
	if len(factions) < 1 {
		log.Println("Can't find any factions for user2")
	} else {
		for _, f := range factions {
			log.Println("Found faction", f, "for user2")
		}
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
	g := &Game{db: db, Gid: gid}
	if g.Select() {
		log.Println("Got", g)
		log.Println("Planets:", g.Planets())
	} else {
		log.Println("game select failed")
	}
}

func TestFourth(t *testing.T) {
	str := []byte("(305,5)")
	log.Println(str)
	p := hexagon.Coord{}
	(&p).Scan(str)
	log.Println("POINT:", p)
}
