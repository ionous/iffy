package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref/class"
	"github.com/ionous/iffy/rt"
	"strings"
)

// IsEmpty determines whether the text contains any characters at all.
type IsEmpty struct {
	Text rt.TextEval
}

func (op *IsEmpty) GetBool(run rt.Runtime) (ret bool, err error) {
	if t, e := op.Text.GetText(run); e != nil {
		err = errutil.New("IsEmpty.Text", e)
	} else if len(t) == 0 {
		ret = true
	}
	return
}

// Includes determines whether text contains part.
type Includes struct {
	Text rt.TextEval
	Part rt.TextEval
}

func (op *Includes) GetBool(run rt.Runtime) (ret bool, err error) {
	if text, e := op.Text.GetText(run); e != nil {
		err = errutil.New("Includes.Text", e)
	} else if part, e := op.Part.GetText(run); e != nil {
		err = errutil.New("Includes.Part", e)
	} else {
		ret = strings.Contains(text, part)
	}
	return
}

type ClassName struct {
	Obj rt.ObjectEval
}

func (op *ClassName) GetText(run rt.Runtime) (ret string, err error) {
	if obj, e := op.Obj.GetObject(run); e != nil {
		err = errutil.New("ClassName.Obj", e)
	} else {
		ret = class.FriendlyName(obj.GetClass())
	}
	return
}
