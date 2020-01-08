package ephemera

// Queue provides a wrapper which can write to a db.... or not.
type Queue interface {
	Prep(which string, keys ...Col)
	// for now, panics on error
	Write(which string, args ...interface{}) Queued
}

// Col describes a column in Queue.
type Col struct {
	Name, Type string
}

// Queued provides an opaque return value for rows written by Queues
type Queued struct {
	id int64
}
