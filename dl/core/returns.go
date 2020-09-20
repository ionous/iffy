package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

// Returns allows a local variable to be used as an output.
type Returns struct {
	Name  string
	Using *Activity
}

func (*Returns) Compose() composer.Spec {
	return composer.Spec{
		Name:  "returns",
		Spec:  "Return the variable {name:text} {?using}",
		Group: "variables",
		Desc:  "Return: Return the value of the named variable computed during using.",
	}
}

const returnNotImplemented = errutil.Error("return not implemented")

func (op *Returns) GetBool(run rt.Runtime) (ret bool, err error) {
	err = returnNotImplemented
	return
}

func (op *Returns) GetNumber(run rt.Runtime) (ret float64, err error) {
	err = returnNotImplemented
	return
}

func (op *Returns) GetText(run rt.Runtime) (ret string, err error) {
	err = returnNotImplemented
	return
}

func (op *Returns) GetNumberStream(run rt.Runtime) (ret rt.Iterator, err error) {
	err = returnNotImplemented
	return
}

func (op *Returns) GetTextStream(run rt.Runtime) (ret rt.Iterator, err error) {
	err = returnNotImplemented
	return
}

// func (op *Returns) GetBool(run rt.Runtime) (ret bool, err error) {
// 	err = op.run(run, func(p interface{}) (err error) {
// 		ret, err = ExpandBool(run, p)
// 		return
// 	})
// 	return
// }

// func (op *Returns) run(run rt.Runtime, cb func(interface{}) error) (err error) {
// 	k := returnScope{name: op.Name}
// 	run.PushScope(&k)
// 	if e := rt.RunOne(run, op.Using); e != nil {
// 		err = e
// 	} else {
// 		err = cb(k.v)
// 	}
// 	run.PopScope()
// 	return
// }

// needs more thought with re: both new Value and pattern Prologues
// type returnScope struct {
// 	name string
// 	v    interface{}
// }

// func (k *returnScope) GetVariable(n string) (interface{}, error) {
// 	return k.v, nil
// }

// // note: the command SetVar helps ensure that "v" is a primitive type.
// func (k *returnScope) SetVariable(n string, v interface{}) (err error) {
// 	if n == k.name {
// 		k.v = v
// 	} else {
// 		err = rt.UnknownVariable(n)
// 	}
// 	return
// }
