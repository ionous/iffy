package core

import (
	"bytes"
	"strings"

	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

// IsEmpty determines whether the text contains any characters at all.
type IsEmpty struct {
	Text rt.TextEval
}

// Includes determines whether text contains part.
type Includes struct {
	Text rt.TextEval
	Part rt.TextEval
}

// Join combines multiple pieces of text.
type Join struct {
	Sep   rt.TextEval
	Parts []rt.TextEval // fix? this should probably? be text list eveal
}

func (*IsEmpty) Compose() composer.Spec {
	return composer.Spec{
		Name:  "is_empty",
		Group: "strings",
		Desc:  "Is Empty: True if the text is empty.",
	}
}

func (op *IsEmpty) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else {
		b := len(t.String()) == 0
		ret = g.BoolOf(b)
	}
	return
}

func (*Includes) Compose() composer.Spec {
	return composer.Spec{
		Name:  "includes",
		Group: "strings",
		Desc:  "Includes Text: True if text contains text.",
	}
}

func (op *Includes) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if text, e := safe.GetText(run, op.Text); e != nil {
		err = cmdErrorCtx(op, "Text", e)
	} else if part, e := safe.GetText(run, op.Part); e != nil {
		err = cmdErrorCtx(op, "Part", e)
	} else {
		contains := strings.Contains(text.String(), part.String())
		ret = g.BoolOf(contains)
	}
	return
}

func (*Join) Compose() composer.Spec {
	return composer.Spec{
		Name:  "join",
		Group: "strings",
		Desc:  "Join Strings: Returns multiple pieces of text as a single new piece of text.",
	}
}

func (op *Join) GetText(run rt.Runtime) (ret g.Value, err error) {
	if sep, e := safe.GetOptionalText(run, op.Sep, ""); e != nil {
		err = cmdErrorCtx(op, "Sep", e)
	} else {
		var buf bytes.Buffer
		sep := sep.String()
		for _, txt := range op.Parts {
			if str, e := safe.GetText(run, txt); e != nil {
				err = cmdErrorCtx(op, "Part", e)
				break
			} else {
				if buf.Len() > 0 {
					buf.WriteString(sep)
				}
				buf.WriteString(str.String())
			}
		}
		if err == nil {
			str := buf.String()
			ret = g.StringOf(str)
		}
	}
	return
}

// if sep, e := safe.GetOptionalText(run, op.Sep, ""); e != nil {
// 	err = e
// } else if it, e := safe.GetTextList(run, op.Parts); e != nil {
// 	err = e
// } else {
// 	var buf bytes.Buffer
// 	for it.HasNext() {
// 		var txt string
// 		if e := it.GetNext(&txt); e != nil {
// 			err = e
// 			break
// 		} else {
// 			if buf.Len() > 0 {
// 				buf.WriteString(sep)
// 			}
// 			buf.WriteString(txt)
// 		}
// 	}
// 	if err == nil {
// 		ret = buf.String()
// 	}
// }
// return
