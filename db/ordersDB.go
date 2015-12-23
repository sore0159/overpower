package db

import (
	"fmt"
	"mule/mydb"
	"mule/overpower"
)

func (d DB) ValidOrder(gid, fid, source, target, size int) (ok bool) {
	planets, ok := d.GetPidPlanets(gid, source, target)
	if !ok {
		return false
	}
	if len(planets) != 2 {
		return false
	}
	if size < 1 {
		return true
	}
	var sTotal int
	if planets[0].Pid() == source {
		sTotal = planets[0].Parts()
	} else if planets[1].Pid() == source {
		sTotal = planets[1].Parts()
	} else {
		return false
	}
	orders, ok := d.GetAllSourceOrders(gid, source)
	if !ok {
		return false
	}
	for _, o := range orders {
		if o.Target() == target {
			if curSize := o.Size(); curSize >= size {
				return true
			} else {
				sTotal -= (size - curSize)
			}
		} else {
			sTotal -= o.Size()
		}
		if sTotal < 0 {
			return false
		}
	}
	return true
}

func (d DB) SetOrder(gid, fid, source, target, size int) (ok bool) {
	o := NewOrder()
	o.gid, o.fid, o.source, o.target, o.size = gid, fid, source, target, size
	return mydb.Upsert(d.db, o)
}

func (d DB) DropAllGidOrders(gid int) (ok bool) {
	query := fmt.Sprintf("DELETE FROM orders WHERE gid = %d", gid)
	return mydb.ExecIf(d.db, query)
}

func (d DB) GetAllSourceOrders(gid, pid int) (orders []overpower.Order, ok bool) {
	query := fmt.Sprintf("SELECT %s FROM orders WHERE gid = %d AND source = %d", ODSQLVAL, gid, pid)
	return d.GetOrdersQuery(query)
}

func (d DB) GetAllGidOrders(gid int) (orders []overpower.Order, ok bool) {
	query := fmt.Sprintf("SELECT %s FROM orders WHERE gid = %d", ODSQLVAL, gid)
	return d.GetOrdersQuery(query)
}

func (d DB) GetOrdersQuery(query string) (orders []overpower.Order, ok bool) {
	ords := []*Order{}
	maker := func() mydb.Rower {
		o := NewOrder()
		return o
	}
	if !mydb.Get(d.db, query, &ords, maker) {
		return nil, false
	}
	orders = make([]overpower.Order, len(ords))
	for i, o := range ords {
		orders[i] = o
	}
	return orders, true
}
