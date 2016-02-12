package db

import (
	"mule/overpower"
)

func (d DB) GetLaunchRecord(conditions ...interface{}) (overpower.LaunchRecord, error) {
	list, err := d.GetLaunchRecords(conditions...)
	if my, bad := Check(err, "get LaunchRecord failure"); bad {
		return nil, my
	}
	if err = IsUnique(len(list)); err != nil {
		return nil, err
	}
	return list[0], nil
}

// Example conditions:  C{"gid",1} C{"owner","mule"}, nil
func (d DB) GetLaunchRecords(conditions ...interface{}) ([]overpower.LaunchRecord, error) {
	group := NewLaunchRecordGroup()
	err := d.getGroup(group, conditions)
	if my, bad := Check(err, "get LaunchRecords failure", "conditions", conditions); bad {
		return nil, my
	}
	list := group.List
	converted := convertLaunchRecords2OP(list...)
	return converted, nil
}
