package core

import (
	"regexp"

	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

type Matches struct {
	Text    rt.TextEval
	Pattern string
	// fix: should transform into a different command probably during compile
	exp *regexp.Regexp `if:"internal"`
	err error
}

// Compose defines a spec for the composer editor.
func (*Matches) Compose() composer.Spec {
	return composer.Spec{
		Group: "matching",
		Desc:  `Matches: Determine whether the specified text is similar to the specified regular expression.`,
		Spec:  "{text:text_eval} matches {pattern:text}",
	}
}

func (op *Matches) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if text, e := safe.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else if exp, e := op.getRegexp(); e != nil {
		err = cmdError(op, e)
	} else {
		b := exp.MatchString(text.String())
		ret = g.BoolOf(b)
	}
	return
}

func (op *Matches) getRegexp() (ret *regexp.Regexp, err error) {
	if e := op.err; e != nil {
		err = e
	} else if exp := op.exp; exp != nil {
		ret = exp
	} else if exp, e := regexp.Compile(op.Pattern); e != nil {
		op.err = err
		err = e
	} else {
		op.exp = exp
		ret = exp
	}
	return
}
