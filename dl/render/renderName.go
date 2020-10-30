package render

import (
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
)

// Name handles changing a template like {.boombip} into text.
// if the name is a variable containing an object name: return the printed object name.
// if the name is a variable with some other text: return that text.
// if the name isn't a variable but refers to some object: return that object's printed object name.
// otherwise, its an error.
type Name struct {
	Name string
}

func (op *Name) Compose() composer.Spec {
	return composer.Spec{
		Name:  "render_name",
		Group: "internal",
	}
}
func (op *Name) GetText(run rt.Runtime) (ret string, err error) {
	// uppercase names are assumed to be requests for object names.
	if name := op.Name; lang.IsCapitalized(name) {
		ret, err = op.getPrintedNamedOf(run, name)
	} else {
		// first check if there's a variable of the requested name
		switch v, e := run.GetField(object.Variables, op.Name); e.(type) {
		default:
			err = cmdError(op, e)
		case rt.UnknownTarget, rt.UnknownField:
			// if there was no such variable, then it's probably an object name
			ret, err = op.getPrintedNamedOf(run, name)
		case nil:
			// get the text from the variable
			if n, e := v.GetText(); e != nil {
				err = cmdError(op, e)
			} else if strings.HasPrefix(n, "#") {
				// if its an object id, get its printed name
				ret, err = op.getPrintedNamedOf(run, n)
			} else {
				// if its not, just assume the author was asking for the variable's text
				ret = n
			}
		}
	}
	return
}

func (op *Name) getPrintedNamedOf(run rt.Runtime, objectName string) (ret string, err error) {
	if printedName, e := rt.GetText(run, &core.Buffer{core.NewActivity(
		&pattern.DetermineAct{
			"printName",
			core.NewArgs(&core.FromVar{
				Name:  &core.Text{objectName},
				Flags: 0}),
		})}); e != nil {
		err = cmdError(op, e)
	} else {
		ret = printedName
	}
	return
}

func cmdError(op composer.Slat, e error) error {
	return errutil.Append(&core.CommandError{Cmd: op}, e)
}
