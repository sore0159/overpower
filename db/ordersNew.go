package db

type Order struct {
	gid    int
	fid    int
	source int
	target int
	size   int
}

func NewOrder() *Order {
	return &Order{
	//
	}
}

func (o *Order) Gid() int {
	return o.gid
}
func (o *Order) Fid() int {
	return o.fid
}
func (o *Order) Source() int {
	return o.source
}
func (o *Order) Target() int {
	return o.target
}
func (o *Order) Size() int {
	return o.size
}
func (o *Order) SetSize(size int) {
	o.size = size
}
