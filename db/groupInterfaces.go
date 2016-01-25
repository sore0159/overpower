package db

import (
	"mule/mydb"
)

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

func (d DB) updateGroup(group UpdateGroup) error {
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
}

func (d DB) getGroup(group SelectGroup, conditions []interface{}) error {
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
}
func (d DB) makeGroup(group InsertGroup) error {
	sqlers := group.InsertList()
	if len(sqlers) == 0 {
		return nil
	}
	table := group.SQLTable()
	cols := group.InsertCols()
	scanCols := group.InsertScanCols()
	query := mydb.InsertQ(table, cols, scanCols)
	return mydb.Insert(d.db(), query, cols, scanCols, sqlers...)
}

func (d DB) dropItems(table string, conditions []interface{}) error {
	query, args, err := mydb.DeleteQA(table, conditions)
	if my, bad := Check(err, "dropitem failure", "table", table, "conditions", conditions); bad {
		return my
	}
	return d.mustExec(query, args...)
}

func (d DB) dropItemsIf(table string, conditions []interface{}) error {
	query, args, err := mydb.DeleteQA(table, conditions)
	if my, bad := Check(err, "dropitem failure", "table", table, "conditions", conditions); bad {
		return my
	}
	_, err = d.db().Exec(query, args...)
	return err
}
