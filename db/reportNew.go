package db

import "mule/mydb"

type Report struct {
	gid      int
	fid      int
	turn     int
	contents mydb.StringList
}

func NewReport() *Report {
	return &Report{
		contents: *mydb.NewStringList(),
	}
}

func (r *Report) AddContent(x string) {
	r.contents = append(r.contents, x)
}
func (r *Report) Contents() []string {
	return []string(r.contents)
}
func (r *Report) Gid() int {
	return r.gid
}
func (r *Report) Fid() int {
	return r.fid
}
func (r *Report) Turn() int {
	return r.turn
}

func (item *Report) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return item.gid
	case "fid":
		return item.fid
	case "turn":
		return item.turn
	case "contents":
		return item.contents
	}
	return nil
}

func (item *Report) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &item.gid
	case "fid":
		return &item.fid
	case "turn":
		return &item.turn
	case "contents":
		return &item.contents
	}
	return nil
}

func (item *Report) SQLTable() string {
	return "reports"
}

func (group *ReportGroup) SQLTable() string {
	return "reports"
}

func (group *ReportGroup) SelectCols() []string {
	return []string{
		"gid",
		"fid",
		"turn",
		"contents",
	}
}

func (group *ReportGroup) UpdateCols() []string {
	return nil
}

func (group *ReportGroup) PKCols() []string {
	return []string{
		"gid",
		"fid",
		"turn",
		"contents",
	}
}

func (group *ReportGroup) InsertCols() []string {
	return []string{
		"gid",
		"fid",
		"turn",
		"contents",
	}
}

func (group *ReportGroup) InsertScanCols() []string {
	return nil
}
