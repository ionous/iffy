package next

import (
	"strconv"

	"github.com/divan/num2words"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
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
		Desc:  "Num as text: Writes a number using numerals, eg. '1'.",
		Group: "printing",
	}
}

func (p *PrintNum) GetText(run rt.Runtime) (ret string, err error) {
	if n, e := rt.GetNumber(run, p.Num); e != nil {
		err = e
	} else if s := strconv.FormatFloat(n, 'g', -1, 64); len(s) > 0 {
		ret = s
	} else {
		ret = "<num>"
	}
	return
}

// Compose defines a spec for the composer editor.
func (*PrintNumWord) Compose() composer.Spec {
	return composer.Spec{
		Name:  "print_num_word",
		Desc:  "Num in words: Writes a number in plain english: eg. 'one'",
		Group: "printing",
	}
}

func (p *PrintNumWord) GetText(run rt.Runtime) (ret string, err error) {
	if n, e := rt.GetNumber(run, p.Num); e != nil {
		err = e
	} else if s := num2words.Convert(int(n)); len(s) > 0 {
		ret = s
	} else {
		ret = "<num>"
	}
	return
}
