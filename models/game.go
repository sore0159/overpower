package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"mule/mydb/db"
	gp "mule/mydb/group"
	sq "mule/mydb/sql"
	"mule/overpower"
)

type Game struct {
	GID       int            `json:"gid"`
	Owner     string         `json:"owner"`
	Name      string         `json:"name"`
	Turn      int            `json:"turn"`
	Autoturn  int            `json:"-"`
	FreeAutos int            `json:"freeautos"`
	Password  sql.NullString `json:"-"`
	ToWin     int            `json:"towin"`
	HighScore int            `json:"highscore"`
	sql       gp.SQLStruct
}

// --------- BEGIN GENERIC METHODS ------------ //

func NewGame() *Game {
	return &Game{
	//
	}
}

type GameIntf struct {
	item *Game
}

func (item *Game) Intf() overpower.GameDat {
	return &GameIntf{item}
}

func (i GameIntf) DELETE() {
	i.item.sql.DELETE = true
}

func (item *Game) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return item.GID
	case "owner":
		return item.Owner
	case "name":
		return item.Name
	case "turn":
		return item.Turn
	case "autoturn":
		return item.Autoturn
	case "freeautos":
		return item.FreeAutos
	case "password":
		return item.Password
	case "towin":
		return item.ToWin
	case "highscore":
		return item.HighScore
	}
	return nil
}

func (item *Game) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &item.GID
	case "owner":
		return &item.Owner
	case "name":
		return &item.Name
	case "turn":
		return &item.Turn
	case "autoturn":
		return &item.Autoturn
	case "freeautos":
		return &item.FreeAutos
	case "password":
		return &item.Password
	case "towin":
		return &item.ToWin
	case "highscore":
		return &item.HighScore
	}
	return nil
}
func (item *Game) SQLTable() string {
	return "game"
}

func (i GameIntf) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		*Game
		HasPassword bool    `json:"haspassword"`
		AutoDays    [7]bool `json:"autodays"`
	}{
		Game:        i.item,
		HasPassword: i.HasPassword(),
		AutoDays:    i.AutoDays(),
	})
}
func (i GameIntf) UnmarshalJSON(data []byte) error {
	i.item = &Game{}
	return json.Unmarshal(data, i.item)
}

func (i GameIntf) GID() int {
	return i.item.GID
}

func (i GameIntf) Owner() string {
	return i.item.Owner
}

func (i GameIntf) SetOwner(x string) {
	if i.item.Owner == x {
		return
	}
	i.item.Owner = x
	i.item.sql.UPDATE = true
}

func (i GameIntf) Name() string {
	return i.item.Name
}

func (i GameIntf) SetName(x string) {
	if i.item.Name == x {
		return
	}
	i.item.Name = x
	i.item.sql.UPDATE = true
}

func (i GameIntf) Turn() int {
	return i.item.Turn
}

func (i GameIntf) IncTurn() {
	i.item.Turn += 1
	i.item.sql.UPDATE = true
}
func (i GameIntf) SetTurn(x int) {
	if i.item.Turn == x {
		return
	}
	i.item.Turn = x
	i.item.sql.UPDATE = true
}

func (i GameIntf) AutoDays() (days [7]bool) {
	sum := i.item.Autoturn
	for j := 0; j < 7; j++ {
		if sum%2 == 1 {
			days[j] = true
		}
		sum = sum / 2
	}
	return
}

func (i GameIntf) SetAutoDays(days [7]bool) {
	var sum int
	for j, b := range days {
		if b {
			sum += 1 << uint32(j)
		}
	}
	i.SetAutoturn(sum)
}

func (i GameIntf) Autoturn() int {
	return i.item.Autoturn
}

func (i GameIntf) SetAutoturn(x int) {
	if i.item.Autoturn == x {
		return
	}
	i.item.Autoturn = x
	i.item.sql.UPDATE = true
}

func (i GameIntf) FreeAutos() int {
	return i.item.FreeAutos
}

func (i GameIntf) SetFreeAutos(x int) {
	if i.item.FreeAutos == x {
		return
	}
	i.item.FreeAutos = x
	i.item.sql.UPDATE = true
}

func (i GameIntf) HasPassword() bool {
	return i.item.Password.Valid && i.item.Password.String != ""
}
func (i GameIntf) IsPassword(x string) bool {
	return i.item.Password.String == x
}

func (i GameIntf) SetPassword(x sql.NullString) {
	if i.item.Password == x {
		return
	}
	i.item.Password = x
	i.item.sql.UPDATE = true
}

func (i GameIntf) ToWin() int {
	return i.item.ToWin
}

func (i GameIntf) SetToWin(x int) {
	if i.item.ToWin == x {
		return
	}
	i.item.ToWin = x
	i.item.sql.UPDATE = true
}

func (i GameIntf) HighScore() int {
	return i.item.HighScore
}

func (i GameIntf) SetHighScore(x int) {
	if i.item.HighScore == x {
		return
	}
	i.item.HighScore = x
	i.item.sql.UPDATE = true
}

// --------- END GENERIC METHODS ------------ //
// --------- BEGIN CUSTOM METHODS ------------ //

// --------- END CUSTOM METHODS ------------ //
// --------- BEGIN GROUP ------------ //

type GameGroup struct {
	List []*Game
}

func NewGameGroup() *GameGroup {
	return &GameGroup{
		List: []*Game{},
	}
}

func (item *Game) SQLGroup() gp.SQLGrouper {
	return NewGameGroup()
}

func (group *GameGroup) New() gp.SQLer {
	item := NewGame()
	group.List = append(group.List, item)
	return item
}

func (group *GameGroup) UpdateList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.UPDATE && !item.sql.INSERT && !item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *GameGroup) InsertList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.INSERT && !item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *GameGroup) DeleteList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *GameGroup) SQLTable() string {
	return "game"
}

func (group *GameGroup) PKCols() []string {
	return []string{
		"gid",
	}
}

func (group *GameGroup) InsertCols() []string {
	return []string{
		"owner",
		"name",
		"turn",
		"autoturn",
		"freeautos",
		"password",
		"towin",
		"highscore",
	}
}

func (group *GameGroup) InsertScanCols() []string {
	return []string{
		"gid",
	}
}

func (group *GameGroup) SelectCols() []string {
	return []string{
		"gid",
		"owner",
		"name",
		"turn",
		"autoturn",
		"freeautos",
		"password",
		"towin",
		"highscore",
	}
}

func (group *GameGroup) UpdateCols() []string {
	return []string{
		"owner",
		"name",
		"turn",
		"autoturn",
		"freeautos",
		"password",
		"towin",
		"highscore",
	}
}

// --------- END GROUP ------------ //
// --------- BEGIN SESSION ------------ //
type GameSession struct {
	*GameGroup
	*gp.Session
}

func NewGameSession(d db.DBer) *GameSession {
	group := NewGameGroup()
	return &GameSession{
		GameGroup: group,
		Session:   gp.NewSession(group, d),
	}
}

func (s *GameSession) Select(conditions ...interface{}) ([]overpower.GameDat, error) {
	cur := len(s.GameGroup.List)
	err := s.Session.Select(conditions...)
	if my, bad := Check(err, "Game select failed", "conditions", conditions); bad {
		return nil, my
	}
	return convertGame2Intf(s.GameGroup.List[cur:]...), nil
}

func (s *GameSession) SelectWhere(where sq.Condition) ([]overpower.GameDat, error) {
	cur := len(s.GameGroup.List)
	err := s.Session.SelectWhere(where)
	if my, bad := Check(err, "Game SelectWhere failed", "where", where); bad {
		return nil, my
	}
	return convertGame2Intf(s.GameGroup.List[cur:]...), nil
}

// --------- END SESSION  ------------ //
// --------- BEGIN UTILS ------------ //

func convertGame2Struct(list ...overpower.GameDat) ([]*Game, error) {
	mylist := make([]*Game, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(GameIntf); ok {
			mylist = append(mylist, t.item)
		} else {
			return nil, errors.New("bad Game struct type for conversion")
		}
	}
	return mylist, nil
}

func convertGame2Intf(list ...*Game) []overpower.GameDat {
	converted := make([]overpower.GameDat, len(list))
	for i, item := range list {
		converted[i] = item.Intf()
	}
	return converted
}

func GameTableCreate(d db.DBer) error {
	query := `create table game(
	gid SERIAL PRIMARY KEY,
	owner varchar(20) NOT NULL UNIQUE,
	name varchar(20) NOT NULL,
	turn int NOT NULL DEFAULT 0,
	autoturn int NOT NULL DEFAULT 0,
	freeautos int NOT NULL DEFAULT 0,
	towin int NOT NULL,
	highscore int NOT NULL DEFAULT 0,
	winner text DEFAULT NULL,
	password varchar(20) DEFAULT NULL
);`
	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed Game table creation", "query", query); bad {
		return my
	}
	return nil
}

func GameTableDelete(d db.DBer) error {
	query := "DROP TABLE IF EXISTS game CASCADE"
	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed Game table deletion", "query", query); bad {
		return my
	}
	return nil
}

// --------- END UTILS ------------ //
