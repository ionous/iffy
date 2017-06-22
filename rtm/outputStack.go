package rtm

import (
	"github.com/ionous/errutil"
	"io"
)

// OutputStack provides a stack of io.Writer(s)
// Iffy uses these to print and consolidate bits of output text.
type OutputStack struct {
	output []io.Writer
}

// PushWriter to activate the passed writer.
func (os *OutputStack) PushWriter(w io.Writer) {
	os.output = append(os.output, w)
}

// PushWriter to deactivate the most recently pushed writer.
func (os *OutputStack) PopWriter() {
	os.output = os.output[:len(os.output)-1]
}

// Write implements io.Writer, delegating to the active writer.
func (os *OutputStack) Write(p []byte) (ret int, err error) {
	if cnt := len(os.output); cnt == 0 {
		err = errutil.New("output stack empty")
	} else {
		w := os.output[len(os.output)-1]
		ret, err = w.Write(p)
	}
	return
}
