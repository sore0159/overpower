package db

import (
	"mule/hexagon"
)

type LaunchRecord struct {
	gid    int
	fid    int
	turn   int
	source hexagon.Coord
	target hexagon.Coord
	size   int
}

func NewLaunchRecord() *LaunchRecord {
	return &LaunchRecord{
	//
	}
}

func (o *LaunchRecord) Gid() int {
	return o.gid
}
func (o *LaunchRecord) Fid() int {
	return o.fid
}
func (o *LaunchRecord) Turn() int {
	return o.turn
}
func (o *LaunchRecord) Source() hexagon.Coord {
	return o.source
}
func (o *LaunchRecord) Target() hexagon.Coord {
	return o.target
}
func (o *LaunchRecord) Size() int {
	return o.size
}
func (o *LaunchRecord) SetSize(size int) {
	o.size = size
}
func (item *LaunchRecord) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return item.gid
	case "fid":
		return item.fid
	case "turn":
		return item.turn
	case "sourcex":
		return item.source[0]
	case "sourcey":
		return item.source[1]
	case "targetx":
		return item.target[0]
	case "targety":
		return item.target[1]
	case "size":
		return item.size
	}
	return nil
}

func (item *LaunchRecord) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &item.gid
	case "fid":
		return &item.fid
	case "turn":
		return &item.turn
	case "sourcex":
		return &item.source[0]
	case "sourcey":
		return &item.source[1]
	case "targetx":
		return &item.target[0]
	case "targety":
		return &item.target[1]
	case "size":
		return &item.size
	}
	return nil
}

func (item *LaunchRecord) SQLTable() string {
	return "launchrecords"
}

func (group *LaunchRecordGroup) SQLTable() string {
	return "launchrecords"
}

func (group *LaunchRecordGroup) SelectCols() []string {
	return []string{
		"gid",
		"fid",
		"turn",
		"sourcex",
		"sourcey",
		"targetx",
		"targety",
		"size",
	}
}

func (group *LaunchRecordGroup) UpdateCols() []string {
	return nil
}

func (group *LaunchRecordGroup) PKCols() []string {
	return []string{
		"gid",
		"fid",
		"turn",
		"sourcex",
		"sourcey",
		"targetx",
		"targety",
	}
}

func (group *LaunchRecordGroup) InsertCols() []string {
	return []string{
		"gid",
		"fid",
		"turn",
		"sourcex",
		"sourcey",
		"targetx",
		"targety",
		"size",
	}
}

func (group *LaunchRecordGroup) InsertScanCols() []string {
	return nil
}
