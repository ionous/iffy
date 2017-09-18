package tests

import (
	"bytes"
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"github.com/kr/pretty"
)

type Bool struct{ rt.BoolEval }
type Number struct{ rt.NumberEval }
type Text struct{ rt.TextEval }
type Execute struct{ rt.Execute }

func (n Execute) MatchLine(run rt.Runtime, expected string) (err error) {
	var buf bytes.Buffer
	if e := rt.WritersBlock(run, &buf, func() error {
		return n.Execute.Execute(run)
	}); e != nil {
		err = e
	} else {
		got := buf.String()
		if d := pretty.Diff(got, expected); len(d) > 0 {
			err = errutil.New(d, got, expected)
		}
	}
	return
}

func (n Execute) GetLines(run rt.Runtime) (ret []string, err error) {
	var lines printer.Lines
	if e := rt.WritersBlock(run, &lines, func() error {
		return n.Execute.Execute(run)
	}); e != nil {
		err = e
	} else {
		ret = lines.Lines()
	}
	return
}
func (n Execute) MatchLines(run rt.Runtime, expected []string) (err error) {
	if got, e := n.GetLines(run); e != nil {
		err = e
	} else if d := pretty.Diff(got, expected); len(d) > 0 {
		err = errutil.New(d, got, expected)
	}
	return
}

func (n Bool) Match(run rt.Runtime, expected bool) (err error) {
	if got, e := n.GetBool(run); e != nil {
		err = e
	} else if got != expected {
		err = errutil.Fmt("%v != %v (got != expected)", got, expected)
	}
	return
}

func (n Number) Match(run rt.Runtime, expected float64) (err error) {
	if got, e := n.GetNumber(run); e != nil {
		err = e
	} else if got != expected {
		err = errutil.Fmt("%v != %v (got != expected)", got, expected)
	}
	return
}

func (n Text) Match(run rt.Runtime, expected string) (err error) {
	if got, e := n.GetText(run); e != nil {
		err = e
	} else if got != expected {
		err = errutil.Fmt("%v != %v (got != expected)", got, expected)
	}
	return
}
