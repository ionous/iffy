package scope

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

// Parameters maps names to arbitrary evaluations.
type Parameters map[string]interface{}

// GetVariable returns the value at 'name', the caller is responsible for determining the type.
func (p Parameters) GetVariable(name string) (ret interface{}, err error) {
	if i, ok := p[name]; !ok {
		err = errutil.New("variable", name, "not found")
	} else {
		ret = i
	}
	return
}

// SetVariable writes the value of 'v' into the value at 'name'.
func (p Parameters) SetVariable(name string, v interface{}) (err error) {
	p[name] = v // FIX: any sort of validation? ex. ensure the value is baked.
	return
}

// Bake distills the parameters into particular values.
func (p Parameters) Bake(run rt.Runtime) (ret Parameters, err error) {
	out := make(Parameters)
	for k, i := range p {
		if v, e := bake(run, i); e != nil {
			err = e
			break
		} else {
			out[k] = v
		}
	}
	if err == nil {
		ret = out
	}
	return
}

func bake(run rt.Runtime, i interface{}) (ret interface{}, err error) {
	switch eval := i.(type) {
	default:
		err = errutil.Fmt("unknown type %T", i)
	case rt.BoolEval:
		if v, e := eval.GetBool(run); e != nil {
			err = e
		} else {
			ret = v
		}
	case rt.NumberEval:
		if v, e := eval.GetNumber(run); e != nil {
			err = e
		} else {
			ret = v
		}
	case rt.TextEval:
		if v, e := eval.GetText(run); e != nil {
			err = e
		} else {
			ret = v
		}
	case rt.NumListEval:
		if vals, e := rt.GetNumList(run, eval); e != nil {
			err = e
		} else {
			ret = vals
		}
	case rt.TextListEval:
		if vals, e := rt.GetTextList(run, eval); e != nil {
			err = e
		} else {
			ret = vals
		}
	}
	return
}
