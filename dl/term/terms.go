package term

import (
	"github.com/ionous/iffy/affine"
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
	run    rt.Runtime
	values termValues
}

func MakeTerms(run rt.Runtime) Terms {
	return Terms{run: run} // delay creating the map
}

type termValues map[string]*termValue

type termValue struct {
	affinity affine.Affinity // alt: if needed might point back to the original term structure
	value    rt.Value
}

func (ps *Terms) addTerm(field string, affinity affine.Affinity, value rt.Value) {
	if ps.values == nil {
		ps.values = make(termValues)
	}
	ps.values[field] = &termValue{affinity, value}
}

// GetField returns the value at 'name', the caller is responsible for determining the type.
func (ps *Terms) GetField(target, field string) (ret rt.Value, err error) {
	if target != object.Variables {
		err = rt.UnknownTarget{target}
	} else if p, ok := ps.values[field]; !ok {
		err = rt.UnknownField{target, field}
	} else {
		ret, err = generic.CopyValue(p.affinity, p.value)
	}
	return
}

// SetField writes (a copy of) the passed value into the term at 'name'.
func (ps *Terms) SetField(target, field string, val rt.Value) (err error) {
	if target != object.Variables {
		err = rt.UnknownTarget{target}
	} else if p, ok := ps.values[field]; !ok {
		err = rt.UnknownField{target, field}
	} else if v, e := generic.CopyValue(p.affinity, val); e != nil {
		err = e
	} else {
		p.value = v
	}
	return
}
