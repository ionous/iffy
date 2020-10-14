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
	run    rt.Runtime
	values parameterValues
}

func MakeTerms(run rt.Runtime) Terms {
	return Terms{run: run} // delay creating the map
}

type parameterValues map[string]*generic.Value

func (ps *Terms) write(name string, typeName string, value interface{}) {
	if ps.values == nil {
		ps.values = make(parameterValues)
	}
	ps.values[name] = generic.NewValue(typeName, value)
}

// GetVariable returns the value at 'name', the caller is responsible for determining the type.
func (ps *Terms) GetField(target, field string) (ret rt.Value, err error) {
	if target != object.Variables {
		err = rt.UnknownTarget{target}
	} else if v, ok := ps.values[field]; !ok {
		err = rt.UnknownField{target, field}
	} else {
		ret = v
	}
	return
}

// SetVariable writes the value of 'v' into the value at 'name'.
func (ps *Terms) SetField(target, field string, val rt.Value) (err error) {
	if target != object.Variables {
		err = rt.UnknownTarget{target}
	} else if p, ok := ps.values[field]; !ok {
		err = rt.UnknownField{target, field}
	} else {
		err = p.SetValue(ps.run, val)
	}
	return
}
