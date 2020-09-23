package pattern

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/tables"
)

type NumParam struct {
	Name string
}
type BoolParam struct {
	Name string
}
type TextParam struct {
	Name string
}
type ObjectParam struct {
	Name, Kind string
}
type NumListParam struct {
	Name string
}
type TextListParam struct {
	Name string
}

// parameters in the future might have defaults...
// something similar might be used for local variables.
// we could also -- in some far future land -- code generate things.
type Parameter interface {
	String() string
	Prepare(Parameters)
}

func (n *NumParam) String() string {
	return n.Name
}
func (n *NumParam) Prepare(p Parameters) {
	var v float64
	p.values[n.Name] = generic.NewValue(tables.PRIM_DIGI, v)
}

func (n *BoolParam) String() string {
	return n.Name
}
func (n *BoolParam) Prepare(p Parameters) {
	var v bool
	p.values[n.Name] = generic.NewValue(tables.PRIM_BOOL, v)
}

func (n *TextParam) String() string {
	return n.Name
}
func (n *TextParam) Prepare(p Parameters) {
	var v string
	p.values[n.Name] = generic.NewValue(tables.PRIM_TEXT, v)
}

func (n *ObjectParam) String() string {
	return n.Name
}
func (n *ObjectParam) Prepare(p Parameters) {
	var v string
	p.values[n.Name] = generic.NewValue(tables.PRIM_TEXT, v)
}

func (n *NumListParam) String() string {
	return n.Name
}
func (n *NumListParam) Prepare(p Parameters) {
	var v []float64
	p.values[n.Name] = generic.NewValue("num_list", v)
}

func (n *TextListParam) String() string {
	return n.Name
}
func (n *TextListParam) Prepare(p Parameters) {
	var v []string
	p.values[n.Name] = generic.NewValue("text_list", v)
}

// Parameters implements a Scope mapping names to specified parameters.
// The only current user is pattern.FromPattern::Stitch()
// It stores values from indexed and key name arguments ( originally specified as evals. )
// Its pushed into scope so the names can be used as a source of values for rt.Runtime::GetVariable().
// ( ex. For use with the commands GetVar{},  SimpleNoun{}, ProperNoun{}, ObjectName{}, ... )
type Parameters struct {
	run    rt.Runtime
	values parameterValues
}
type parameterValues map[string]*generic.Value

// GetVariable returns the value at 'name', the caller is responsible for determining the type.
func (ps *Parameters) GetVariable(name string) (ret rt.Value, err error) {
	if i, ok := ps.values[name]; !ok {
		err = rt.UnknownVariable(name)
	} else {
		ret = i
	}
	return
}

// SetVariable writes the value of 'v' into the value at 'name'.
func (ps *Parameters) SetVariable(name string, val rt.Value) (err error) {
	// FIX: any sort of validation? ex. ensure the value is baked ( ie. some sort of primitive or slice of primitives. )
	if p, ok := ps.values[name]; !ok {
		err = rt.UnknownVariable(name)
	} else {
		err = p.SetValue(ps.run, val)
	}
	return
}
