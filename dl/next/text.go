package next

import (
	"bytes"
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

// Join combines multiple pieces of text.
type Join struct {
	Elems rt.TextListEval
	Sep   rt.TextEval
}

func (op *Join) GetText(run rt.Runtime) (ret string, err error) {
	if sep, e := rt.GetOptionalText(run, op.Sep); e != nil {
		err = e
	} else if it, e := rt.GetTextStream(run, op.Elems); e != nil {
		err = e
	} else {
		var buf bytes.Buffer
		for it.HasNext() {
			if txt, e := it.GetText(); e != nil {
				err = e
				break
			} else {
				if buf.Len() > 0 {
					buf.WriteString(sep)
				}
				buf.WriteString(txt)
			}
		}
		if err == nil {
			ret = buf.String()
		}
	}
	return
}
