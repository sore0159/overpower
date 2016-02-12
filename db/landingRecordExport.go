package db

import (
	"mule/overpower"
)

func (d DB) GetLandingRecord(conditions ...interface{}) (overpower.LandingRecord, error) {
	list, err := d.GetLandingRecords(conditions...)
	if my, bad := Check(err, "get LandingRecord failure"); bad {
		return nil, my
	}
	if err = IsUnique(len(list)); err != nil {
		return nil, err
	}
	return list[0], nil
}

// Example conditions:  C{"gid",1} C{"owner","mule"}, nil
func (d DB) GetLandingRecords(conditions ...interface{}) ([]overpower.LandingRecord, error) {
	group := NewLandingRecordGroup()
	err := d.getGroup(group, conditions)
	if my, bad := Check(err, "get LandingRecords failure", "conditions", conditions); bad {
		return nil, my
	}
	list := group.List
	converted := convertLandingRecords2OP(list...)
	return converted, nil
}
