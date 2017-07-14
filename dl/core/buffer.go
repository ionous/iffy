package core

import (
	"github.com/ionous/iffy/rt"
)

// Buffer collects text said by other statements via a SpanPrinter, and returns it as a string.
type Buffer struct {
	Buffer []rt.Execute
}

func (buf *Buffer) GetText(run rt.Runtime) (ret string, err error) {
	var span SpanPrinter
	run.PushWriter(&span)
	if e := rt.ExecuteList(run, buf.Buffer); e != nil {
		err = e
	} else {
		ret = span.String()
	}
	run.PopWriter()
	return
}
