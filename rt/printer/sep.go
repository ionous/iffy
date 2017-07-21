package printer

import (
	"io"
)

// Sep implements io.Writer, treating every Write a new word.
type Sep struct {
	io.Writer
	mid, last string // seperators
	queue     string // last string sent to Write()
	cnt       int    // number of non-zero writes to the underlying writer.
}

// AndSeparator creates a phrase: a, b, c, and d.
// Note: spacing between words is left to printer.Span.
func AndSeparator(w io.Writer) *Sep {
	return &Sep{Writer: w, mid: ",", last: ", and"}
}

// OrSeparator creates a phrase: a, b, c, or d.
// Note: spacing between words is left to printer.Span.
func OrSeparator(w io.Writer) *Sep {
	return &Sep{Writer: w, mid: ",", last: ", or"}
}

// Write implements io.Writer, spacing
func (l *Sep) Write(p []byte) (ret int, err error) {
	if len(p) > 0 {
		s := string(p)
		err = l.flush(l.mid)
		l.queue = s
		ret = len(s)
	}
	return
}

// Close writes all current lines, with appropriate separators, to the passed output.
// It does not Close the wrapped stream.
func (l *Sep) Close() (err error) {
	return l.flush(l.last)
}

// Flush empties the queue
func (l *Sep) flush(sep string) (err error) {
	if len(l.queue) > 0 {
		if l.cnt != 0 {
			_, e := io.WriteString(l.Writer, sep)
			err = e
		}
		if err == nil {
			_, e := io.WriteString(l.Writer, l.queue)
			err = e
		}
		l.queue = ""
		l.cnt++
	}
	return
}
