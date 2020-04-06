package core

import (
	"github.com/ionous/iffy/assign"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
)

// Determine runes a pattern
type Determine struct {
	Name       string
	Parameters scope.Parameters // serialized directly, b/c its unique per invocation.
}

func (*Determine) Compose() composer.Spec {
	return composer.Spec{
		Name:  "set_num_list_var",
		Group: "internal",
		Desc:  "Set Number List Variable: Sets the named variable to the passed number list.",
	}
}

func (det *Determine) Execute(run rt.Runtime) (err error) {
	if p, e := run.GetField(det.Name, object.Pattern); e != nil {
		err = e
	} else if exec, ok := p.(rt.Execute); !ok {
		err = assign.Mismatch(det.Name, exec, p)
	} else {
		err = scopeBlock(run, det.Parameters, func() (err error) {
			err = exec.Execute(run)
			return
		})
	}
	return
}

func (det *Determine) GetBool(run rt.Runtime) (ret bool, err error) {
	if p, e := run.GetField(det.Name, object.Pattern); e != nil {
		err = e
	} else if eval, ok := p.(rt.BoolEval); !ok {
		err = assign.Mismatch(det.Name, eval, p)
	} else {
		err = scopeBlock(run, det.Parameters, func() (err error) {
			ret, err = eval.GetBool(run)
			return
		})
	}
	return
}

func (det *Determine) GetNumber(run rt.Runtime) (ret float64, err error) {
	if p, e := run.GetField(det.Name, object.Pattern); e != nil {
		err = e
	} else if eval, ok := p.(rt.NumberEval); !ok {
		err = assign.Mismatch(det.Name, eval, p)
	} else {
		err = scopeBlock(run, det.Parameters, func() (err error) {
			ret, err = eval.GetNumber(run)
			return
		})
	}
	return
}

func (det *Determine) GetText(run rt.Runtime) (ret string, err error) {
	if p, e := run.GetField(det.Name, object.Pattern); e != nil {
		err = e
	} else if eval, ok := p.(rt.TextEval); !ok {
		err = assign.Mismatch(det.Name, eval, p)
	} else {
		err = scopeBlock(run, det.Parameters, func() (err error) {
			ret, err = eval.GetText(run)
			return
		})
	}
	return
}

func (det *Determine) GetNumberStream(run rt.Runtime) (ret rt.Iterator, err error) {
	if p, e := run.GetField(det.Name, object.Pattern); e != nil {
		err = e
	} else if eval, ok := p.(rt.NumListEval); !ok {
		err = assign.Mismatch(det.Name, eval, p)
	} else {
		err = scopeBlock(run, det.Parameters, func() (err error) {
			ret, err = eval.GetNumberStream(run)
			return
		})
	}
	return
}

func (det *Determine) GetTextStream(run rt.Runtime) (ret rt.Iterator, err error) {
	if p, e := run.GetField(det.Name, object.Pattern); e != nil {
		err = e
	} else if eval, ok := p.(rt.TextListEval); !ok {
		err = assign.Mismatch(det.Name, eval, p)
	} else {
		err = scopeBlock(run, det.Parameters, func() (err error) {
			ret, err = eval.GetTextStream(run)
			return
		})
	}
	return
}

func scopeBlock(run rt.Runtime, params scope.Parameters, fn func() error) (err error) {
	if scope, e := params.Bake(run); e != nil {
		err = e
	} else {
		run.PushScope(scope)
		err = fn()
		run.PopScope()
	}
	return
}
