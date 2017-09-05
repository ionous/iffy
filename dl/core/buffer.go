package core

import (
	"bytes"
	"github.com/ionous/iffy/rt"
)

// Buffer collects text said by other statements and returns them as a string.
// Unlike PrintSpan, it does not add or alter spaces between writes.
type Buffer struct {
	Buffer rt.ExecuteList
}

func (buf *Buffer) GetText(run rt.Runtime) (ret string, err error) {
	var span bytes.Buffer
	if e := buf.Buffer.Execute(rt.Writer(run, &span)); e != nil {
		err = e
	} else {
		ret = span.String()
	}
	return
}
