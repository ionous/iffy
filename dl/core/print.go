package core

import (
	"github.com/divan/num2words"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"io"
	"strconv"
)

// PrintSpan writes text inline, with spaces between words.
type PrintSpan struct {
	Block rt.ExecuteList
}

// PrintBracket sandwiches text inside parenthesis.
type PrintBracket struct {
	Block rt.ExecuteList
}

// PrintSlash separates text with a left-leaning slash.
type PrintSlash struct {
	Block rt.ExecuteList
}

// PrintList writes words separated with commas, ending with an "and".
type PrintList struct {
	Block rt.ExecuteList
}

// PrintNum writes a number using numerals, eg. "1".
type PrintNum struct {
	Num rt.NumberEval
}

// PrintNumWord writes a number using english: eg. "one".
type PrintNumWord struct {
	Num rt.NumberEval
}

// Say writes a piece of text.
type Say struct {
	Text rt.TextEval
}

func (p *PrintSpan) Execute(run rt.Runtime) error {
	return rt.WritersBlock(run, printer.Spanning(run.Writer()), func() error {
		return p.Block.Execute(run)
	})
}

func (p *PrintBracket) Execute(run rt.Runtime) error {
	return rt.WritersBlock(run, printer.Bracket(run.Writer()), func() error {
		return p.Block.Execute(run)
	})
}

func (p *PrintList) Execute(run rt.Runtime) error {
	return rt.WritersBlock(run, printer.AndSeparator(run.Writer()), func() error {
		return p.Block.Execute(run)
	})
}

func (p *PrintNum) Execute(run rt.Runtime) (err error) {
	if n, e := p.Num.GetNumber(run); e != nil {
		err = e
	} else if s := strconv.FormatFloat(n, 'g', -1, 64); len(s) > 0 {
		_, err = io.WriteString(run, s)
	} else {
		_, err = io.WriteString(run, "<num>")
	}
	return err
}

func (p *PrintNumWord) Execute(run rt.Runtime) (err error) {
	if n, e := p.Num.GetNumber(run); e != nil {
		err = e
	} else if s := num2words.Convert(int(n)); len(s) > 0 {
		_, err = io.WriteString(run, s)
	} else {
		_, err = io.WriteString(run, "<num>")
	}
	return err
}

func (p *PrintSlash) Execute(run rt.Runtime) error {
	return rt.WritersBlock(run, printer.Slash(run.Writer()), func() error {
		return p.Block.Execute(run)
	})
}

func (p *Say) Execute(run rt.Runtime) (err error) {
	if s, e := p.Text.GetText(run); e != nil {
		err = e
	} else if len(s) > 0 {
		_, err = io.WriteString(run, s)
	}
	return
}

func Print(run rt.Runtime, text rt.TextEval) (err error) {
	if s, e := text.GetText(run); e != nil {
		err = e
	} else {
		_, err = io.WriteString(run, s)
	}
	return err
}
