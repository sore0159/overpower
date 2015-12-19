package overpower

import "mule/mydb"

type Game interface {
	Turn() int
	IncTurn()
	Gid() int
	Name() string
	Owner() string
	HasPW() bool
	IsPwd(string) bool
	//
	mydb.Updater
}

type Faction interface {
	Gid() int
	Fid() int
	Owner() string
	Name() string
	Done() bool
	SetDone(bool)
	//
	mydb.Updater
}
