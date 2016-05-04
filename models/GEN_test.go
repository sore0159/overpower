package models

import (
	"fmt"
	"log"
	//	"mule/hexagon"
	"mule/mybad"
	"mule/overpower"
	"testing"
)

func TestOne(t *testing.T) {
	log.Println("TEST ONE")
}

func TestTwo(t *testing.T) {
	db, err := LoadDB()
	ErrCheck(err)
	log.Println("Loaded DB")
	logE, failE := db.Transact(DropTest)
	ErrCheck(failE)
	if logE != nil {
		log.Println("LOG ERROR:")
		log.Println(logE.(*mybad.MuleError).MuleError())
	}

	logE, failE = db.Transact(MakeTest)
	ErrCheck(failE)
	if logE != nil {
		log.Println("LOG ERROR:")
		log.Println(logE.(*mybad.MuleError).MuleError())
	}
	log.Println("MakeTest complete!")
}

func DropTest(m *Manager) (logE, failE error) {
	gList, err := m.Game().Select("owner", "Testing_User")
	if my, bad := Check(err, "Droptest search error"); bad {
		return nil, my
	}
	if len(gList) == 0 {
		log.Println("No previous test games found")
		return nil, nil
	}
	for _, g := range gList {
		log.Println("DELETING PREVIOUS TESTING GAME", g.GID())
		g.DELETE()
	}
	return nil, m.Close()
}

func MakeTest(m *Manager) (logE, failE error) {
	g := &Game{
		Owner: "Testing_User",
		Name:  "TestGame",
		ToWin: 10,
	}
	m.CreateGame(g)
	err := m.Close()
	if my, bad := Check(err, "creategame testing failure", "game", g); bad {
		return nil, my
	}
	log.Println("MADE GAME", g)
	madeF := make([]*Faction, 5)
	for i := 0; i < 4; i++ {
		f := &Faction{
			GID:   g.GID,
			Owner: fmt.Sprintf("Tester%d", i),
			Name:  fmt.Sprintf("Faction%d", i),
		}
		madeF[i] = f
		m.CreateFaction(f)
	}
	f := &Faction{
		GID:   g.GID,
		Owner: "Test",
		Name:  "MainTestFaction",
	}
	madeF[4] = f
	m.CreateFaction(f)
	err = m.Close()
	if my, bad := Check(err, "createfaction testing failure", "factions", madeF); bad {
		return nil, my
	}
	log.Println("MADE FACTIONS", madeF)
	source := NewSource(m, g.GID)
	return nil, overpower.MakeGalaxy(source, false)
}
