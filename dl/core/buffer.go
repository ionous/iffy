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

func (b *Buffer) GetText(run rt.Runtime) (ret string, err error) {
	var buf bytes.Buffer
	if e := rt.WritersBlock(run, &buf, func() error {
		return b.Buffer.Execute(run)
	}); e != nil {
		err = e
	} else {
		ret = buf.String()
	}
	return
}
