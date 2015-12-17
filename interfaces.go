package overpower

type Game interface {
	Turn() int
	IncTurn()
	Gid() int
	Name() string
	Owner() string
	IsPwd(string) bool
}
