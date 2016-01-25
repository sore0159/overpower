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
func (item *Order) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return item.gid
	case "fid":
		return item.fid
	case "source":
		return item.source
	case "target":
		return item.target
	case "size":
		return item.size
	}
	return nil
}

func (item *Order) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &item.gid
	case "fid":
		return &item.fid
	case "source":
		return &item.source
	case "target":
		return &item.target
	case "size":
		return &item.size
	}
	return nil
}

func (item *Order) SQLTable() string {
	return "orders"
}

func (group *OrderGroup) SQLTable() string {
	return "orders"
}

func (group *OrderGroup) SelectCols() []string {
	return []string{
		"gid",
		"fid",
		"source",
		"target",
		"size",
	}
}

func (group *OrderGroup) UpdateCols() []string {
	return []string{
		"size",
	}
}

func (group *OrderGroup) PKCols() []string {
	return []string{
		"gid",
		"fid",
		"source",
		"target",
	}
}

func (group *OrderGroup) InsertCols() []string {
	return []string{
		"gid",
		"fid",
		"source",
		"target",
		"size",
	}
}

func (group *OrderGroup) InsertScanCols() []string {
	return nil
}
