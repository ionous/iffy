package core

import (
	"bytes"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
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

func (op *IsEmpty) GetBool(run rt.Runtime) (ret bool, err error) {
	if t, e := rt.GetText(run, op.Text); e != nil {
		err = errutil.New("IsEmpty.Text", e)
	} else if len(t) == 0 {
		ret = true
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

func (*Join) Compose() composer.Spec {
	return composer.Spec{
		Name:  "join",
		Group: "strings",
		Desc:  "Join Strings: Returns multiple pieces of text as a single new piece of text.",
	}
}

func (op *Join) GetText(run rt.Runtime) (ret string, err error) {
	if sep, e := rt.GetOptionalText(run, op.Sep, ""); e != nil {
		err = e
	} else {
		var buf bytes.Buffer
		for _, txt := range op.Parts {
			if str, e := rt.GetText(run, txt); e != nil {
				err = e
				break
			} else {
				if buf.Len() > 0 {
					buf.WriteString(sep)
				}
				buf.WriteString(str)
			}
		}
		if err == nil {
			ret = buf.String()
		}
	}
	return
}

// if sep, e := rt.GetOptionalText(run, op.Sep, ""); e != nil {
// 	err = e
// } else if it, e := rt.GetTextList(run, op.Parts); e != nil {
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
