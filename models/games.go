package models

type Game struct {
	ID    int
	Name  string
	Owner string
}

func NewGame() *Game {
	return &Game{}
}

func GetGame(id int) *Game {
	g := NewGame()
	g.ID = id
	query := "SELECT name, owner FROM games where id == $1"
	err := DB.QueryRow(query, id).Scan(&(g.Name), &(g.Owner))
	if err != nil {
		Log(err)
		return nil
	}
	return g
}

func MakeGame(name, owner string) *Game {
	g := NewGame()
	g.Name = name
	g.Owner = owner
	query := "INSERT INTO games (name, owner) VALUES($1, $2) RETURNING id"
	err := DB.QueryRow(query, name, owner).Scan(&(g.ID))
	if err != nil {
		Log(err)
		return nil
	}
	return g
}

func (g *Game) Destroy() {
	query := "DELETE FROM games where id = $1"
	res, err := DB.Exec(query, g.ID)
	if err != nil {
		Log("failed to delete game", g.ID, ":", err)
		return
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		Log("failed to delete game", g.ID, ": 0 rows affected")
		return
	}
}
