package next

import (
	"io"
	"strconv"

	"github.com/divan/num2words"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

// PrintNum writes a number using numerals, eg. "1".
type PrintNum struct {
	Num rt.NumberEval
}

// Compose defines a spec for the composer editor.
func (*PrintNum) Compose() composer.Spec {
	return composer.Spec{
		Desc:  "Num as text: Writes a number using numerals, eg. '1'.",
		Group: "printing",
	}
}

func (p *PrintNum) WriteText(run rt.Runtime, w io.Writer) (err error) {
	if n, e := p.Num.GetNumber(run); e != nil {
		err = e
	} else if s := strconv.FormatFloat(n, 'g', -1, 64); len(s) > 0 {
		_, err = io.WriteString(w, s)
	} else {
		_, err = io.WriteString(w, "<num>")
	}
	return
}

// PrintNumWord writes a number using english: eg. "one".
type PrintNumWord struct {
	Num rt.NumberEval
}

// Compose defines a spec for the composer editor.
func (*PrintNumWord) Compose() composer.Spec {
	return composer.Spec{
		Desc:  "Num in words: Writes a number in plain english: eg. 'one'",
		Group: "printing",
	}
}

func (p *PrintNumWord) WriteText(run rt.Runtime, w io.Writer) (err error) {
	if n, e := p.Num.GetNumber(run); e != nil {
		err = e
	} else if s := num2words.Convert(int(n)); len(s) > 0 {
		_, err = io.WriteString(w, s)
	} else {
		_, err = io.WriteString(w, "<num>")
	}
	return
}
