package print

import (
	"bytes"

	"github.com/ionous/iffy/rt/writer"
)

// Sep implements writer.Output, treating every Write as a new word.
type Sep struct {
	out     writer.Output
	last    string       // separators
	pending bytes.Buffer // last string sent to Write()
	cnt     int          // number of non-zero writes to the underlying writer.
}

// AndSeparator creates a phrase: a, b, c, and d.
// Note: spacing between words is left to print.Spacing.
func AndSeparator(w writer.Output) writer.OutputCloser {
	return &Sep{out: w, last: "and"}
}

// OrSeparator creates a phrase: a, b, c, or d.
// Note: spacing between words is left to print.Spacing.
func OrSeparator(w writer.Output) writer.OutputCloser {
	return &Sep{out: w, last: "or"}
}

func (l *Sep) Write(p []byte) (int, error) {
	return l.write(Chunk{p})
}
func (l *Sep) WriteByte(c byte) error {
	_, e := l.write(Chunk{c})
	return e
}
func (l *Sep) WriteRune(r rune) (int, error) {
	return l.write(Chunk{r})
}
func (l *Sep) WriteString(s string) (int, error) {
	return l.write(Chunk{s})
}

// Close writes all pending lines with appropriate separators.
func (l *Sep) Close() error {
	if l.cnt > 1 {
		l.out.WriteRune(',')
	}
	return l.flush(l.last)
}

// Write implements writer.Output, spacing writes with separators.
func (l *Sep) write(c Chunk) (ret int, err error) {
	if !c.IsEmpty() {
		const mid = ","
		if e := l.flush(mid); e != nil {
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
			_, e := l.out.WriteString(sep)
			err = e
		}
		// write the pending text
		if err == nil {
			_, e := l.out.Write(l.pending.Bytes())
			err = e
		}
		l.pending.Reset()
		l.cnt++
	}
	return
}
