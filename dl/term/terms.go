package term

import (
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
)

// Terms implements a Scope mapping names to specified parameters.
// The only current user is pattern.FromPattern::Stitch()
// It stores values from indexed and key name arguments ( originally specified as evals. )
// Its pushed into scope so the names can be used as a source of values for rt.Runtime::GetField().
// ( ex. For use with the commands GetVar{},  SimpleNoun{}, ProperNoun{}, ObjectName{}, ... )
type Terms struct {
	values termValues
}

type termValues map[string]*Value

func (ps *Terms) AddTerm(field string, value rt.Value) *Value {
	if ps.values == nil {
		ps.values = make(termValues)
	}
	v := &Value{value}
	ps.values[field] = v
	return v
}

// GetField returns the value at 'name', the caller is responsible for determining the type.
func (ps *Terms) GetField(target, field string) (ret rt.Value, err error) {
	if target != object.Variables {
		err = rt.UnknownTarget{target}
	} else if p, ok := ps.values[field]; !ok {
		err = rt.UnknownField{target, field}
	} else {
		ret, err = generic.CopyValue(p.value.Affinity(), p.value)
	}
	return
}

// SetField writes (a copy of) the passed value into the term at 'name'.
func (ps *Terms) SetField(target, field string, val rt.Value) (err error) {
	if target != object.Variables {
		err = rt.UnknownTarget{target}
	} else if p, ok := ps.values[field]; !ok {
		err = rt.UnknownField{target, field}
	} else if v, e := generic.CopyValue(p.value.Affinity(), val); e != nil {
		err = e
	} else {
		p.value = v
	}
	return
}
