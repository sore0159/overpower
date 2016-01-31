package db

import (
	"mule/mydb"
)

/*
type DeleteGroup interface {
	PKCols() []string
	SQLTable() string
	DeleteList() []mydb.SQLer
}

type SelectGroup interface {
	SQLTable() string
	SelectCols() []string
	New() mydb.SQLer
}
type UpdateGroup interface {
	SQLTable() string
	UpdateCols() []string
	PKCols() []string
	UpdateList() []mydb.SQLer
}
type InsertGroup interface {
	SQLTable() string
	InsertCols() []string
	InsertScanCols() []string
	InsertList() []mydb.SQLer
}
*/

func (d DB) updateGroup(group mydb.UpdateGrouper) error {
	return mydb.UpdateGroup(d.db(), group)
	/*
		sqlers := group.UpdateList()
		if len(sqlers) == 0 {
			return nil
		}
		table := group.SQLTable()
		cols := group.UpdateCols()
		condCols := group.PKCols()
		query := mydb.UpdateQ(table, cols, condCols)
		allCols := append(cols, condCols...)
		return mydb.Update(d.db(), true, query, allCols, sqlers...)
	*/
}

func (d DB) getGroup(group mydb.SelectGrouper, conditions []interface{}) error {
	return mydb.GetGroup(d.db(), group, conditions)
	/*
		table := group.SQLTable()
		cols := group.SelectCols()
		query, args, err := mydb.SelectQA(table, cols, conditions)
		if my, bad := Check(err, "get interface failure", "table", table, "conditions", conditions, "cols", cols); bad {
			return my
		}
		err = mydb.Get(d.db(), group, query, args...)
		if my, bad := Check(err, "get interface failure", "query", query, "args", args); bad {
			return my
		}
		return nil
	*/
}
func (d DB) makeGroup(group mydb.InsertGrouper) error {
	return mydb.MakeGroup(d.db(), group)
	/*
		sqlers := group.InsertList()
		if len(sqlers) == 0 {
			return nil
		}
		table := group.SQLTable()
		cols := group.InsertCols()
		scanCols := group.InsertScanCols()
		query := mydb.InsertQ(table, cols, scanCols)
		return mydb.Insert(d.db(), query, cols, scanCols, sqlers...)
	*/
}

func (d DB) dropItems(table string, conditions []interface{}) error {
	return mydb.DropItems(d.db(), table, conditions)
	/*
		query, args, err := mydb.DeleteQA(table, conditions)
		if my, bad := Check(err, "dropitem failure", "table", table, "conditions", conditions); bad {
			return my
		}
		return d.mustExec(query, args...)
	*/
}

func (d DB) dropItemsIf(table string, conditions []interface{}) error {
	return mydb.DropItemsIf(d.db(), table, conditions)
	/*
		query, args, err := mydb.DeleteQA(table, conditions)
		if my, bad := Check(err, "dropitem failure", "table", table, "conditions", conditions); bad {
			return my
		}
		_, err = d.db().Exec(query, args...)
		return err
	*/
}

func (d DB) dropGroup(group mydb.DeleteGrouper) error {
	return mydb.DropGroup(d.db(), group)
}

func (d DB) updateItem(table string, set, conditions []interface{}) error {
	return mydb.UpdateItem(d.db(), table, set, conditions)
	/*
		setCols, setArgs, err := set.Split()
		if my, bad := Check(err, "update item failure on set splic", "set", set, "conditions", conditions); bad {
			return my
		}
		condCols, condArgs, err := conditions.Split()
		if my, bad := Check(err, "update item failure on conditions split", "set", set, "conditions", conditions); bad {
			return my
		}
		query := mydb.UpdateQ(table, setCols, condCols)
		args := append(setArgs, condArgs...)
		err = d.mustExec(query, args...)
		if my, bad := Check(err, "update item failure", "table", table, "query", query, "args", args); bad {
			return my
		}
		return nil
	*/
}
