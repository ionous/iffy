package ephemera

type Queued struct {
	id int64
}

// wrapper which can write to a db.... or not.
type Queue interface {
	Prep(which string, keys ...string)
	// for now, panics on error
	Write(which string, args ...interface{}) Queued
}
