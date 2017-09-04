package rt

import (
	"io"
)

// Writer returns a new runtime that uses the passed writer for its output.
func Writer(run Runtime, w io.Writer) Runtime {
	return _Writer{run, w}
}

type _Writer struct {
	Runtime
	Writer io.Writer
}

func (l _Writer) Write(p []byte) (int, error) {
	return l.Writer.Write(p)
}
