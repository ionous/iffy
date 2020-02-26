package unique

import (
	r "reflect"
)

type Stack struct {
	parent TypeRegistry
	Types
}

// NewStack creates a registry such that everything in this will also go to parent,
// but not everything in parent will be available via this.
func NewStack(parent TypeRegistry) *Stack {
	return &Stack{
		parent,
		make(Types),
	}
}

// Register chains the call to the parent registry.
func (cs *Stack) RegisterType(rtype r.Type) (err error) {
	if p := cs.parent; p != nil {
		if e := p.RegisterType(rtype); e != nil {
			err = e
		}
	}
	if err == nil {
		err = cs.Types.RegisterType(rtype)
	}
	return
}
