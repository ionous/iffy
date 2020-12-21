package ident

import (
	"strconv"
)

// Counters is a helper to generate semi-unique names for a group.
type Counters map[string]uint64

func (m *Counters) Next(name string) string {
	c := (*m)[name] + 1
	(*m)[name] = c // COUNTER:#
	return name + "_" + strconv.FormatUint(c, 36)
}
