package printer

import (
	"bytes"
)

// Lines implements io.Writer, buffering every Write as a new line.
type Lines struct {
	lines []string
}

// Lines returns all current lines.
// There is no flush. A new line writer can be constructed instead.
func (l *Lines) Lines() []string {
	return l.lines
}

// Write implements io.Writer, treating every Write as a new line.
func (l *Lines) Write(p []byte) (ret int, err error) {
	var buf bytes.Buffer
	ret, err = buf.Write(p)
	l.lines = append(l.lines, buf.String())
	return
}
