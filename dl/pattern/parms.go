package pattern

import "github.com/ionous/iffy/rt/scope"

type NumParam struct {
	Name string
}
type BoolParam struct {
	Name string
}
type TextParam struct {
	Name string
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
	p[n.Name] = v
}

func (n *BoolParam) String() string {
	return n.Name
}
func (n *BoolParam) Prepare(p Parameters) {
	var v bool
	p[n.Name] = v
}

func (n *TextParam) String() string {
	return n.Name
}
func (n *TextParam) Prepare(p Parameters) {
	var v string
	p[n.Name] = v
}

func (n *NumListParam) String() string {
	return n.Name
}
func (n *NumListParam) Prepare(p Parameters) {
	var v []float64
	p[n.Name] = v
}

func (n *TextListParam) String() string {
	return n.Name
}
func (n *TextListParam) Prepare(p Parameters) {
	var v []string
	p[n.Name] = v
}

// Parameters implements a VariableScope mapping names to specified parameters.
// The only current user is pattern.FromPattern::Stitch()
// It stores values from indexed and key name arguments ( originally specified as evals. )
// Its pushed into scope so the names can be used as a source of values for rt.Runtime::GetVariable().
// ( ex. For use with the commands GetVar{},  CommonNoun{}, ProperNoun{}, ObjectName{}, ... )
type Parameters map[string]interface{}

// GetVariable returns the value at 'name', the caller is responsible for determining the type.
func (p Parameters) GetVariable(name string) (ret interface{}, err error) {
	if i, ok := p[name]; !ok {
		err = scope.UnknownVariable(name)
	} else {
		ret = i
	}
	return
}

// SetVariable writes the value of 'v' into the value at 'name'.
func (p Parameters) SetVariable(name string, v interface{}) (err error) {
	// FIX: any sort of validation? ex. ensure the value is baked ( ie. some sort of primitive or slice of primitives. )
	if _, ok := p[name]; !ok {
		err = scope.UnknownVariable(name)
	} else {
		p[name] = v
	}
	return
}
