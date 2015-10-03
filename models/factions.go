package models

type Faction struct {
	FID  int
	GID  int
	User string
	Name string
}

func NewFaction() *Faction {
	return &Faction{
	//
	}
}

func MakeFaction(gID int, user, name string) *Faction {
	f := NewFaction()
	f.GID = gID
	f.User = user
	f.Name = name
	query := "INSERT INTO factions (gid, user, name) VALUES ($1, $2, $3) RETURNING fid"
	err := DB.QueryRow(query, gID, user, name).Scan(&(f.FID))
	if err != nil {
		Log(err)
		return nil
	}
	return f
}

func GetFaction(gID, fID int) *Faction {
	f := NewFaction()
	f.FID, f.GID = fID, gID
	query := "SELECT user, name FROM factions WHERE gid = $1 AND fid = $2"
	err := DB.QueryRow(query, id).Scan(&(g.Name), &(g.Owner))
	if err != nil {
		Log(err)
		return nil
	}
	return f
}
