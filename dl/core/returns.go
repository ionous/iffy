package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
)

// Returns allows a local variable to be used as an output.
type Returns struct {
	Name  string
	Using []rt.Execute
}

type returnScope struct {
	name string
	v    interface{}
}

func (k *returnScope) GetVariable(n string) (interface{}, error) {
	return k.v, nil
}

// note: the command SetVar helps ensure that "v" is a primitive type.
func (k *returnScope) SetVariable(n string, v interface{}) (err error) {
	if n == k.name {
		k.v = v
	} else {
		err = scope.UnknownVariable(n)
	}
	return
}

func (*Returns) Compose() composer.Spec {
	return composer.Spec{
		Name:  "returns",
		Spec:  "Return the variable {name:text} {?using|ghost}",
		Group: "variables",
		Desc:  "Return: Return the value of the named variable computed during using.",
	}
}

func (op *Returns) run(run rt.Runtime, cb func(interface{}) error) (err error) {
	k := returnScope{name: op.Name}
	run.PushScope(&k)
	if e := rt.RunAll(run, op.Using); e != nil {
		err = e
	} else {
		err = cb(k.v)
	}
	run.PopScope()
	return
}

func (op *Returns) GetBool(run rt.Runtime) (ret bool, err error) {
	err = op.run(run, func(p interface{}) (err error) {
		ret, err = GetBool(run, p)
		return
	})
	return
}

func (op *Returns) GetNumber(run rt.Runtime) (ret float64, err error) {
	err = op.run(run, func(p interface{}) (err error) {
		ret, err = GetNumber(run, p)
		return
	})
	return
}

func (op *Returns) GetText(run rt.Runtime) (ret string, err error) {
	err = op.run(run, func(p interface{}) (err error) {
		ret, err = GetText(run, p)
		return
	})
	return
}

func (op *Returns) GetNumberStream(run rt.Runtime) (ret rt.Iterator, err error) {
	err = op.run(run, func(p interface{}) (err error) {
		ret, err = GetNumbers(run, p)
		return
	})
	return
}

func (op *Returns) GetTextStream(run rt.Runtime) (ret rt.Iterator, err error) {
	err = op.run(run, func(p interface{}) (err error) {
		ret, err = GetTexts(run, p)
		return
	})
	return
}
