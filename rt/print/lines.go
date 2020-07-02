package print

import (
	"bytes"

	"github.com/ionous/iffy/rt/writer"
)

// Lines implements io.Writer, buffering every Write as a new line.
// use MakeChunks to construct a valid line writer.
type Lines struct {
	writer.ChunkOutput
	lines []string
}

func NewLines() *Lines {
	lines := new(Lines)
	writer.InitChunks(lines)
	return lines
}

// Lines returns all current lines.
// There is no flush. A new line writer can be constructed instead.
func (l *Lines) Lines() []string {
	return l.lines
}

// Write implements writer.Output, spacing writes with separators.
func (l *Lines) WriteChunk(c writer.Chunk) (int, error) {
	var buf bytes.Buffer
	n, e := c.WriteTo(&buf)
	l.lines = append(l.lines, buf.String())
	return n, e
}
