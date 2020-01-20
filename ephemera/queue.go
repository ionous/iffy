package ephemera

import "database/sql/driver"

// Queue provides a wrapper which can write to a db.... or not.
type Queue interface {
	Prep(which string, keys ...Col)
	// for now, panics on error
	Write(which string, args ...interface{}) (Queued, error)
}

// Queued provides a semi-opaque return value for rows written by Queues
type Queued struct {
	id int64
}

// Scan converts a database value into a Queued entry. ( opposite of Value )
func (ns *Queued) Scan(value interface{}) (err error) {
	if v, e := driver.DefaultParameterConverter.ConvertValue(value); e != nil {
		err = e
	} else {
		ns.id = v.(int64)
	}
	return
}

// Value converts a Queued entry into a database value. ( opposite of Scan )
func (ns Queued) Value() (driver.Value, error) {
	return ns.id, nil
}
