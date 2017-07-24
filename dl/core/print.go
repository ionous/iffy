package core

import (
	"github.com/divan/num2words"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"io"
	"strconv"
)

type PrintLine struct {
	Block rt.ExecuteList
}

type PrintNum struct {
	Num rt.NumberEval
}

type PrintNumWord struct {
	Num rt.NumberEval
}

type PrintText struct {
	Text rt.TextEval
}

func (p *PrintLine) Execute(run rt.Runtime) (err error) {
	var buf printer.Span
	run.PushWriter(&buf)
	err = p.Block.Execute(run)
	run.PopWriter()
	if err == nil {
		_, err = run.Write(buf.Bytes())
	}
	return
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

func (p *PrintText) Execute(run rt.Runtime) error {
	return Print(run, p.Text)
}

func Print(run rt.Runtime, text rt.TextEval) (err error) {
	if s, e := text.GetText(run); e != nil {
		err = e
	} else {
		_, err = io.WriteString(run, s)
	}
	return err
}
