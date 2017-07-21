package printer

import (
	"github.com/ionous/errutil"
	"io"
)

// Stack provides a stack of io.Writer(s)
// Iffy uses these to print and consolidate bits of output text.
type Stack struct {
	output    []io.Writer
	LastError error
}

// PushWriter to activate the passed writer.
func (os *Stack) PushWriter(w io.Writer) {
	os.output = append(os.output, w)
}

// PopWriter to deactivate the most recently pushed writer.
func (os *Stack) PopWriter() {
	n := len(os.output) - 1
	if top, ok := os.output[n].(io.Closer); ok {
		if e := top.Close(); e != nil {
			os.LastError = e
		}
	}
	os.output = os.output[:n]
}

// Write implements io.Writer, delegating to the active writer.
func (os *Stack) Write(p []byte) (ret int, err error) {
	if cnt := len(os.output); cnt == 0 {
		err = errutil.New("output stack empty")
	} else {
		w := os.output[len(os.output)-1]
		ret, err = w.Write(p)
		if err != nil {
			os.LastError = err
		}
	}
	return
}
