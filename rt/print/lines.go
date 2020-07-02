package print

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

func (l *Lines) Write(p []byte) (int, error) {
	return l.write(Chunk{p})
}
func (l *Lines) WriteByte(c byte) error {
	_, e := l.write(Chunk{c})
	return e
}
func (l *Lines) WriteRune(r rune) (int, error) {
	return l.write(Chunk{r})
}
func (l *Lines) WriteString(s string) (int, error) {
	return l.write(Chunk{s})
}

// Write implements writer.Output, spacing writes with separators.
func (l *Lines) write(c Chunk) (int, error) {
	var buf bytes.Buffer
	n, e := c.WriteTo(&buf)
	l.lines = append(l.lines, buf.String())
	return n, e
}
