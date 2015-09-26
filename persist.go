package main

import (
	"encoding/gob"
	"os"
)

const SAVEDIR = DATADIR + "SV/"

func Save(g *Game, fileName string) error {
	dataFile, err := os.Create(SAVEDIR + fileName)
	if err != nil {
		Log(err)
		return err
	}
	defer dataFile.Close()
	dataEncoder := gob.NewEncoder(dataFile)
	err = dataEncoder.Encode(g)
	if err != nil {
		Log(err)
		return err
	}
	return nil
}

func Load(fileName string) (*Game, error) {
	var g *Game
	dataFile, err := os.Open(SAVEDIR + fileName)
	if err != nil {
		if os.IsNotExist(err) {
			Log("GAME", SAVEDIR+fileName, "NOT FOUND")
			return nil, nil
		}
		Log(err)
		return nil, err
	}
	g = NewGame()
	defer dataFile.Close()
	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(g)
	if err != nil {
		Log(err)
		return nil, err
	}
	return g, nil
}
