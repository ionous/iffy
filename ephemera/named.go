package ephemera

import (
	"database/sql/driver"
	"encoding/json"
)

// opaque row id for name
type Named struct {
	id  int64
	str string
}

func (ns *Named) IsValid() bool {
	return ns.id > 0
}

func (ns *Named) String() string {
	return ns.str
}

func (ns Named) MarshalJSON() ([]byte, error) {
	return json.Marshal(ns.str)
}

// Scan converts a database value into a Named entry. ( opposite of Value )
func (ns *Named) Scan(value interface{}) (err error) {
	if v, e := driver.DefaultParameterConverter.ConvertValue(value); e != nil {
		err = e
	} else {
		ns.id = v.(int64)
	}
	return
}

// Value converts a Named entry into a database value. ( opposite of Scan )
func (ns Named) Value() (driver.Value, error) {
	return ns.id, nil
}
