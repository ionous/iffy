package rtm

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

// NoAncestors provides a default for rtm if one is not set.
type NoAncestors struct{}

func (NoAncestors) GetAncestors(rt.Object) (rt.ObjectStream, error) {
	return NotIt{}, nil
}

// an iterator that always fails
type NotIt struct{}

func (NotIt) HasNext() bool {
	return false
}

func (NotIt) GetNext() (rt.Object, error) {
	return nil, errutil.New("this never has objects")
}
