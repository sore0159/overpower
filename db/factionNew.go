package db

func NewFaction() *Faction {
	return &Faction{}
}

type Faction struct {
	modified bool
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
	f.modified = true
}

func (f *Faction) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return f.gid
	case "fid":
		return f.fid
	case "owner":
		return f.owner
	case "name":
		return f.name
	case "done":
		return f.done
	}
	return nil
}

func (f *Faction) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &f.gid
	case "fid":
		return &f.fid
	case "owner":
		return &f.owner
	case "name":
		return &f.name
	case "done":
		return &f.done
	}
	return nil
}

func (f *Faction) SQLTable() string {
	return "factions"
}

func (group *FactionGroup) SQLTable() string {
	return "factions"
}

func (group *FactionGroup) SelectCols() []string {
	return []string{
		"gid",
		"fid",
		"owner",
		"name",
		"done",
	}
}

func (group *FactionGroup) UpdateCols() []string {
	return []string{
		"done",
	}
}

func (group *FactionGroup) InsertCols() []string {
	return []string{
		"gid",
		"owner",
		"name",
	}
}

func (group *FactionGroup) PKCols() []string {
	return []string{"gid", "fid"}
}

func (group *FactionGroup) InsertScanCols() []string {
	return nil
}
