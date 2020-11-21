package core

import (
	"strconv"

	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

// PrintNum writes a number using numerals, eg. "1".
type PrintNum struct {
	Num rt.NumberEval
}

// PrintNumWord writes a number using english: eg. "one".
type PrintNumWord struct {
	Num rt.NumberEval
}

// Compose defines a spec for the composer editor.
func (*PrintNum) Compose() composer.Spec {
	return composer.Spec{
		Name:  "print_num",
		Spec:  "as text {num:number_eval}",
		Desc:  "A number as text: Writes a number using numerals, eg. '1'.",
		Group: "printing",
	}
}

func (op *PrintNum) GetText(run rt.Runtime) (ret g.Value, err error) {
	if n, e := safe.GetNumber(run, op.Num); e != nil {
		err = cmdError(op, e)
	} else if s := strconv.FormatFloat(n.Float(), 'g', -1, 64); len(s) > 0 {
		ret = g.StringOf(s)
	} else {
		ret = g.StringOf("<num>")
	}
	return
}

// Compose defines a spec for the composer editor.
func (*PrintNumWord) Compose() composer.Spec {
	return composer.Spec{
		Name:  "print_num_word",
		Desc:  "A number in words: Writes a number in plain english: eg. 'one'",
		Group: "printing",
	}
}

func (op *PrintNumWord) GetText(run rt.Runtime) (ret g.Value, err error) {
	if n, e := safe.GetNumber(run, op.Num); e != nil {
		err = cmdError(op, e)
	} else if s, ok := lang.NumToWords(n.Int()); ok {
		ret = g.StringOf(s)
	} else {
		ret = g.StringOf("<num>")
	}
	return
}
