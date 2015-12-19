package db

func NewFaction() *Faction {
	return &Faction{}
}

type Faction struct {
	doneMod bool
	//
	gid   int
	fid   int
	owner string
	name  string
	done  bool
}

func (f *Faction) Gid() int {
	return f.gid
}
func (f *Faction) Fid() int {
	return f.fid
}
func (f *Faction) Owner() string {
	return f.owner
}
func (f *Faction) Name() string {
	return f.name
}
func (f *Faction) Done() bool {
	return f.done
}
func (f *Faction) SetDone(x bool) {
	if f.done == x {
		return
	}
	f.done = x
	f.doneMod = true
}
