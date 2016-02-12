package db

import (
	"database/sql"
	"mule/hexagon"
)

type LandingRecord struct {
	gid               int
	fid               int
	turn              int
	index             int
	size              int
	target            hexagon.Coord
	firstcontroller   sql.NullInt64
	resultcontroller  sql.NullInt64
	resultinhabitants int
}

func NewLandingRecord() *LandingRecord {
	return &LandingRecord{
	//
	}
}

func (mv *LandingRecord) Gid() int {
	return mv.gid
}
func (mv *LandingRecord) Fid() int {
	return mv.fid
}
func (item *LandingRecord) Turn() int {
	return item.turn
}
func (item *LandingRecord) Index() int {
	return item.index
}
func (item *LandingRecord) Size() int {
	return item.size
}
func (item *LandingRecord) Target() hexagon.Coord {
	return item.target
}
func (p *LandingRecord) FirstController() int {
	if p.firstcontroller.Valid {
		return int(p.firstcontroller.Int64)
	}
	return 0
}
func (p *LandingRecord) ResultController() int {
	if p.resultcontroller.Valid {
		return int(p.resultcontroller.Int64)
	}
	return 0
}
func (item *LandingRecord) ResultInhabitants() int {
	return item.resultinhabitants
}
func (mv *LandingRecord) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return mv.gid
	case "fid":
		return mv.fid
	case "turn":
		return mv.turn
	case "index":
		return mv.index
	case "size":
		return mv.size
	case "targetx":
		return mv.target[0]
	case "targety":
		return mv.target[1]
	case "firstcontroller":
		return mv.firstcontroller
	case "resultcontroller":
		return mv.resultcontroller
	case "resultinhabitants":
		return mv.resultinhabitants
	}
	return nil
}
func (mv *LandingRecord) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &mv.gid
	case "fid":
		return &mv.fid
	case "turn":
		return &mv.turn
	case "index":
		return &mv.index
	case "size":
		return &mv.size
	case "targetx":
		return &mv.target[0]
	case "targety":
		return &mv.target[1]
	case "firstcontroller":
		return &mv.firstcontroller
	case "resultcontroller":
		return &mv.resultcontroller
	case "resultinhabitants":
		return &mv.resultinhabitants
	}
	return nil
}

func (mv *LandingRecord) SQLTable() string {
	return "landingrecords"
}

func (group *LandingRecordGroup) SQLTable() string {
	return "landingrecords"
}

func (group *LandingRecordGroup) SelectCols() []string {
	return []string{
		"gid",
		"fid",
		"turn",
		"index",
		"size",
		"targetx",
		"targety",
		"firstcontroller",
		"resultcontroller",
		"resultinhabitants",
	}
}

func (group *LandingRecordGroup) UpdateCols() []string {
	return nil
}

func (group *LandingRecordGroup) PKCols() []string {
	return []string{"gid", "fid", "turn", "index"}
}

func (group *LandingRecordGroup) InsertCols() []string {
	return []string{
		"gid",
		"fid",
		"turn",
		"index",
		"size",
		"targetx",
		"targety",
		"firstcontroller",
		"resultcontroller",
		"resultinhabitants",
	}
}

func (group *LandingRecordGroup) InsertScanCols() []string {
	return nil
}
