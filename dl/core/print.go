package core

import (
	"github.com/ionous/iffy/rt"
	"io"
	"strconv"
)

type PrintNum struct {
	Num rt.NumberEval
}

type PrintText struct {
	Text rt.TextEval
}

type PrintLine struct {
	Block []rt.Execute
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

func (p *PrintText) Execute(run rt.Runtime) (err error) {
	if s, e := p.Text.GetText(run); e != nil {
		err = e
	} else {
		_, err = io.WriteString(run, s)
	}
	return err
}

func (p *PrintLine) Execute(run rt.Runtime) (err error) {
	var buf SpanPrinter
	run.PushWriter(&buf)
	err = rt.ExecuteList(run, p.Block)
	run.PopWriter()
	if err == nil {
		_, err = run.Write(buf.Bytes())
	}
	return
}
