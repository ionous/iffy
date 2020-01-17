package ephemera

// Queue provides a wrapper which can write to a db.... or not.
type Queue interface {
	Prep(which string, keys ...Col)
	// for now, panics on error
	Write(which string, args ...interface{}) (Queued, error)
}

// Col describes a column in Queue.
type Col struct {
	Name, Type, Check string
}

// Queued provides an opaque return value for rows written by Queues
type Queued struct {
	id int64
}

func NamesOf(cols []Col) []string {
	keys := make([]string, 0, len(cols))
	for _, c := range cols {
		if len(c.Name) > 0 {
			keys = append(keys, c.Name)
		}
	}
	return keys
}
