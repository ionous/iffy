package ref

import (
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

type Classes interface {
	// FIX? the rt.Class return is currently needed only for making property parents -- could it be rt.Class instead?
	RegisterClass(r.Type) (rt.Class, error)

	// GetClass compatible with Runtime
	GetClass(string) (rt.Class, bool)
}

type ClassStack struct {
	parent Classes
	ClassMap
}

// NewClasses creates a registry of classes.
// Everything in this new registry will also go to parent
// but not everything in parent will be available via this.
func NewClassStack(parent Classes) *ClassStack {
	return &ClassStack{
		parent,
		make(ClassMap),
	}
}

// RegisterType mimics RegisterClass providing compatiblity with unique.TypeRegistry
func (cs *ClassStack) RegisterType(rtype r.Type) (err error) {
	_, err = cs.RegisterClass(rtype)
	return
}

// RegisterClass chains the call to the parent registry.
func (cs *ClassStack) RegisterClass(rtype r.Type) (ret rt.Class, err error) {
	if cls, e := cs.parent.RegisterClass(rtype); e != nil {
		err = e
	} else {
		cid := id.MakeId(rtype.Name())
		cs.ClassMap[cid] = cls
		ret = cls
	}
	return
}
