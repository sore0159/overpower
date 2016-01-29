package json

import (
	"mule/overpower"
)

func LoadReport(item overpower.Report) *Report {
	return &Report{
		Gid:      item.Gid(),
		Fid:      item.Fid(),
		Turn:     item.Turn(),
		Contents: item.Contents(),
	}
}

func LoadReports(list []overpower.Report) []*Report {
	jList := make([]*Report, len(list))
	for i, item := range list {
		jList[i] = LoadReport(item)
	}
	return jList
}

type Report struct {
	Gid      int      `json:"gid"`
	Fid      int      `json:"fid"`
	Turn     int      `json:"turn"`
	Contents []string `json:"contents"`
}
