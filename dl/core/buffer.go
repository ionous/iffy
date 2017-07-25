package core

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
)

// Buffer collects text said by other statements via a SpanPrinter, and returns it as a string.
type Buffer struct {
	Buffer rt.ExecuteList
}

func (buf *Buffer) GetText(run rt.Runtime) (ret string, err error) {
	var span printer.Span
	if e := buf.Buffer.Execute(rt.Writer(run, &span)); e != nil {
		err = e
	} else {
		ret = span.String()
	}
	return
}
