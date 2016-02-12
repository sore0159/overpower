package db

func NewFaction() *Faction {
	return &Faction{}
}

type Faction struct {
	modified bool
	//
	gid        int
	fid        int
	owner      string
	name       string
	donebuffer int
	score      int
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
func (f *Faction) IsDone() bool {
	return f.donebuffer != 0
}
func (f *Faction) DoneBuffer() int {
	return f.donebuffer
}
func (f *Faction) SetDoneBuffer(x int) {
	if f.donebuffer == x {
		return
	}
	f.donebuffer = x
	f.modified = true
}

func (f *Faction) Score() int {
	return f.score
}
func (f *Faction) SetScore(x int) {
	if f.score == x {
		return
	}
	f.score = x
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
	case "donebuffer":
		return f.donebuffer
	case "score":
		return f.score
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
	case "donebuffer":
		return &f.donebuffer
	case "score":
		return &f.score
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
		"donebuffer",
		"score",
	}
}

func (group *FactionGroup) UpdateCols() []string {
	return []string{
		"donebuffer",
		"score",
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
