package models

import (
	"log"
	"mule/mybad"
	"testing"
)

var RUNUPDATE byte = 0

func TestUpdateTables(t *testing.T) {
	if RUNUPDATE == 1 {
		log.Println("UPDATING TABLES")
		db, err := LoadDB()
		ErrCheck(err)
		err = DropAllTables(db)
		ErrCheck(err)
		log.Println("Tables dropped!")
		err = CreateAllTables(db)
		ErrCheck(err)
		log.Println("Tables created!")
	}
}

func ErrCheck(err error) {
	if err == nil {
		return
	}
	if my, ok := err.(*mybad.MuleError); ok {
		log.Println(my.MuleError())
	}
	panic(err)
}
