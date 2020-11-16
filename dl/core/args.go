package core

import (
	"strconv"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

type Argument struct {
	Name string // argument name
	From Assignment
}

type Arguments struct {
	Args []*Argument
}

func (*Argument) Compose() composer.Spec {
	return composer.Spec{
		Name:  "argument",
		Spec:  "its {name:variable_name} is {from:assignment}",
		Group: "patterns",
	}
}

func (*Arguments) Compose() composer.Spec {
	return composer.Spec{
		Name:  "arguments",
		Spec:  " when {arguments%args+argument}",
		Group: "patterns",
	}
}

//
func (op *Arguments) Distill(run rt.Runtime, out *g.Record) (err error) {
	k := out.Kind()
	for _, arg := range op.Args {
		if name, e := getParamName(k, arg.Name); e != nil {
			err = errutil.Append(err, e)
		} else if val, e := arg.From.GetAssignedValue(run); e != nil {
			err = errutil.Append(err, e)
		} else if e := out.SetNamedField(name, val); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// change a argument name ( which could be an index ) into a valid param name
// fix: this should happen at assembly time...
func getParamName(k *g.Kind, arg string) (ret string, err error) {
	if usesIndex := len(arg) > 1 && arg[:1] == "$"; !usesIndex {
		ret = arg
	} else if storedIdx, e := strconv.Atoi(arg[1:]); e != nil {
		err = errutil.New("couldnt parse index", arg)
	} else if i := storedIdx - 1; i < 0 || i >= k.NumField() {
		err = errutil.New("field", arg, "not found")
	} else {
		ret = k.Field(i).Name
	}
	return
}
