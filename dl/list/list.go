package list

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
)

var Slats = []composer.Slat{
	(*Len)(nil),
	(*At)(nil),
}

func cmdError(op composer.Slat, e error) error {
	return errutil.Append(&core.CommandError{Cmd: op}, e)
}

// slice: chop out a new list
// splice: one or many, changes the original array, returns the removed elements
// push back
// pop back
// push front
// pop front
// concat list
// for each
// sort ( w/ pattern )
