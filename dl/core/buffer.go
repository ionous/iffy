package core

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
)

// Buffer collects text said by other statements via a SpanPrinter, and returns it as a string.
type Buffer struct {
	Buffer []rt.Execute
}

func (buf *Buffer) GetText(run rt.Runtime) (ret string, err error) {
	var span printer.Span
	run.PushWriter(&span)
	if e := rt.ExecuteList(buf.Buffer).Execute(run); e != nil {
		err = e
	} else {
		ret = span.String()
	}
	run.PopWriter()
	return
}
