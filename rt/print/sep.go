package print

import (
	"bytes"

	"github.com/ionous/iffy/rt/writer"
)

// Sep implements writer.Output, treating every Write as a new word.
type Sep struct {
	writer.ChunkOutput
	target    writer.Output
	mid, last string       // separators
	pending   bytes.Buffer // last string sent to Write()
	cnt       int          // number of non-zero writes to the underlying writer.
}

// AndSeparator creates a phrase: a, b, c, and d.
// Note: spacing between words is left to print.Spacing.
func AndSeparator(w writer.Output) writer.OutputCloser {
	sep := &Sep{target: w, mid: ",", last: "and"}
	writer.InitChunks(sep)
	return sep
}

// OrSeparator creates a phrase: a, b, c, or d.
// Note: spacing between words is left to print.Spacing.
func OrSeparator(w writer.Output) writer.OutputCloser {
	sep := &Sep{target: w, mid: ",", last: "or"}
	writer.InitChunks(sep)
	return sep
}

// Close writes all pending lines with appropriate separators.
func (l *Sep) Close() error {
	if l.cnt > 1 {
		l.target.WriteRune(',')
	}
	return l.flush(l.last)
}

// Write implements writer.Output, spacing writes with separators.
func (l *Sep) WriteChunk(c writer.Chunk) (ret int, err error) {
	if !c.IsEmpty() {
		if e := l.flush(l.mid); e != nil {
			err = e
		} else {
			ret, err = c.WriteTo(&l.pending)
		}
	}
	return
}

// Flush writes pending text, prefixed if needed with a separator
func (l *Sep) flush(sep string) (err error) {
	// pending text pending, write it.
	if l.pending.Len() > 0 {
		// separate text already written
		if l.cnt != 0 {
			_, e := l.target.WriteString(sep)
			err = e
		}
		// write the pending text
		if err == nil {
			_, e := l.target.Write(l.pending.Bytes())
			err = e
		}
		l.pending.Reset()
		l.cnt++
	}
	return
}
