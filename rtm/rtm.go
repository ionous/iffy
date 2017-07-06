package rtm

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
)

type Rtm struct {
	*ref.Classes
	*ref.Objects
	*ref.Relations
	ScopeStack
	OutputStack
	Randomizer
	Ancestors
}

func NewRtm(c *ref.Classes, o *ref.Objects, r *ref.Relations) *Rtm {
	rtm := &Rtm{
		Classes:   c,
		Objects:   o,
		Relations: r,
		Ancestors: NoAncestors{},
	}
	rtm.PushScope(scope.ModelFinder(rtm))
	rtm.Randomizer.Reset(1) // FIX: time?
	return rtm
}

// Ancestors is compatible with the rt.Runtime
type Ancestors interface {
	GetAncestors(rt.Object) (rt.ObjectStream, error)
}

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
