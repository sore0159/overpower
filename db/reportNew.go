package db

type Report struct {
	gid      int
	fid      int
	turn     int
	contents []string
}

func NewReport() *Report {
	return &Report{
		contents: []string{},
	}
}

func (r *Report) AddContent(x string) {
	r.contents = append(r.contents, x)
}
func (r *Report) Contents() []string {
	return r.contents
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
