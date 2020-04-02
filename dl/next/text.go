package next

import (
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

// IsEmpty determines whether the text contains any characters at all.
type IsEmpty struct {
	Text rt.TextEval
}

func (op *IsEmpty) GetBool(run rt.Runtime) (ret bool, err error) {
	if t, e := rt.GetText(run, op.Text); e != nil {
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
	if text, e := rt.GetText(run, op.Text); e != nil {
		err = errutil.New("Includes.Text", e)
	} else if part, e := rt.GetText(run, op.Part); e != nil {
		err = errutil.New("Includes.Part", e)
	} else {
		ret = strings.Contains(text, part)
	}
	return
}
