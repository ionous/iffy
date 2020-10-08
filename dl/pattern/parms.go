package pattern

import (
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/tables"
)

type NumParam struct {
	Name string
	Init rt.NumberEval
}
type BoolParam struct {
	Name string
	Init rt.BoolEval
}
type TextParam struct {
	Name string
	Init rt.TextEval
}
type ObjectParam struct {
	Name, Kind string
	Init       rt.TextEval
}
type NumListParam struct {
	Name string
	Init rt.NumListEval
}
type TextListParam struct {
	Name string
	Init rt.TextListEval
}

// parameters in the future might have defaults...
// something similar might be used for local variables.
// we could also -- in some far future land -- code generate things.
type Parameter interface {
	String() string
	Prepare(rt.Runtime, *Parameters) error
}

func (n *NumParam) String() string {
	return n.Name
}

func (n *NumParam) Prepare(run rt.Runtime, p *Parameters) (err error) {
	if v, e := rt.GetOptionalNumber(run, n.Init, 0); e != nil {
		err = e
	} else {
		p.write(n.Name, tables.PRIM_DIGI, v)
	}
	return
}

func (n *BoolParam) String() string {
	return n.Name
}

func (n *BoolParam) Prepare(run rt.Runtime, p *Parameters) (err error) {
	if v, e := rt.GetOptionalBool(run, n.Init, false); e != nil {
		err = e
	} else {
		p.write(n.Name, tables.PRIM_BOOL, v)
	}
	return
}

func (n *TextParam) String() string {
	return n.Name
}

func (n *TextParam) Prepare(run rt.Runtime, p *Parameters) (err error) {
	if v, e := rt.GetOptionalText(run, n.Init, ""); e != nil {
		err = e
	} else {
		p.write(n.Name, tables.PRIM_TEXT, v)
	}
	return
}

func (n *ObjectParam) String() string {
	return n.Name
}

func (n *ObjectParam) Prepare(run rt.Runtime, p *Parameters) (err error) {
	if v, e := rt.GetOptionalText(run, n.Init, ""); e != nil {
		err = e
	} else {
		p.write(n.Name, tables.PRIM_TEXT, v)
	}
	return
}

func (n *NumListParam) String() string {
	return n.Name
}

func (n *NumListParam) Prepare(run rt.Runtime, p *Parameters) (err error) {
	if vs, e := rt.GetOptionalNumbers(run, n.Init, nil); e != nil {
		err = e
	} else {
		p.write(n.Name, "num_list", vs)
	}
	return
}

func (n *TextListParam) String() string {
	return n.Name
}

func (n *TextListParam) Prepare(run rt.Runtime, p *Parameters) (err error) {
	if vs, e := rt.GetOptionalTexts(run, n.Init, nil); e != nil {
		err = e
	} else {
		p.write(n.Name, "text_list", vs)
	}
	return
}

// Parameters implements a Scope mapping names to specified parameters.
// The only current user is pattern.FromPattern::Stitch()
// It stores values from indexed and key name arguments ( originally specified as evals. )
// Its pushed into scope so the names can be used as a source of values for rt.Runtime::GetField().
// ( ex. For use with the commands GetVar{},  SimpleNoun{}, ProperNoun{}, ObjectName{}, ... )
type Parameters struct {
	run    rt.Runtime
	values parameterValues
}
type parameterValues map[string]*generic.Value

func (ps *Parameters) write(name string, typeName string, value interface{}) {
	if ps.values == nil {
		ps.values = make(parameterValues)
	}
	ps.values[name] = generic.NewValue(typeName, value)
}

// GetVariable returns the value at 'name', the caller is responsible for determining the type.
func (ps *Parameters) GetField(target, field string) (ret rt.Value, err error) {
	if target != object.Variables {
		err = rt.UnknownTarget{target}
	} else if i, ok := ps.values[field]; !ok {
		err = rt.UnknownField{target, field}
	} else {
		ret = i
	}
	return
}

// SetVariable writes the value of 'v' into the value at 'name'.
func (ps *Parameters) SetField(target, field string, val rt.Value) (err error) {
	if target != object.Variables {
		err = rt.UnknownTarget{target}
	} else if p, ok := ps.values[field]; !ok {
		err = rt.UnknownField{target, field}
	} else {
		err = p.SetValue(ps.run, val)
	}
	return
}
