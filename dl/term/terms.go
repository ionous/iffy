package term

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
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

func (ps *Terms) AddTerm(field string, value g.Value) *Value {
	if ps.values == nil {
		ps.values = make(termValues)
	}
	v := &Value{value}
	ps.values[field] = v
	return v
}

// GetField returns the value at 'name', the caller is responsible for determining the type.
func (ps *Terms) GetField(target, field string) (ret g.Value, err error) {
	if target != object.Variables {
		err = rt.UnknownTarget{target}
	} else if p, ok := ps.values[field]; !ok {
		err = rt.UnknownField{target, field}
	} else {
		ret = p.value
	}
	return
}

// SetField writes (a copy of) the passed value into the term at 'name'.
func (ps *Terms) SetField(target, field string, val g.Value) (err error) {
	if target != object.Variables {
		err = rt.UnknownTarget{target}
	} else if p, ok := ps.values[field]; !ok {
		err = rt.UnknownField{target, field}
	} else if a := p.value.Affinity(); a != val.Affinity() {
		err = errutil.New("value is not", a)
	} else {
		p.value = val
	}
	return
}
