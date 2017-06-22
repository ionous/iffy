package rtm

import (
	"bytes"
)

// LineWriter implements io.Writer, treating every Write as a new line.
type LineWriter struct {
	lines []string
}

// Lines returns all current lines.
// There is no flush. A new line writer can be constructed instead.
func (l *LineWriter) Lines() []string {
	return l.lines
}

// Write implements io.Writer, treating every Write as a new line.
func (l *LineWriter) Write(p []byte) (ret int, err error) {
	var buf bytes.Buffer
	ret, err = buf.Write(p)
	l.lines = append(l.lines, buf.String())
	return
}
