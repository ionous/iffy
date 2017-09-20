package ident

import (
	"strconv"
)

// Counters is a helper to generate semi-unique names for a group.
type Counters map[Id]int

func (c Counters) NewName(groupName string) string {
	groupId := IdOf(groupName)
	cnt := c[groupId] + 1
	c[groupId] = cnt
	return groupId.Name + "#" + strconv.Itoa(cnt)
}
