package db

import (
	"mule/overpower"
)

func (d DB) GetBattleRecord(conditions ...interface{}) (overpower.BattleRecord, error) {
	list, err := d.GetBattleRecords(conditions...)
	if my, bad := Check(err, "get BattleRecord failure"); bad {
		return nil, my
	}
	if err = IsUnique(len(list)); err != nil {
		return nil, err
	}
	return list[0], nil
}

// Example conditions:  C{"gid",1} C{"owner","mule"}, nil
func (d DB) GetBattleRecords(conditions ...interface{}) ([]overpower.BattleRecord, error) {
	group := NewBattleRecordGroup()
	err := d.getGroup(group, conditions)
	if my, bad := Check(err, "get BattleRecords failure", "conditions", conditions); bad {
		return nil, my
	}
	list := group.List
	converted := convertBattleRecords2OP(list...)
	return converted, nil
}
