package printer

import (
	"bytes"
	"io"
)

// Sep implements io.Writer, treating every Write a new word.
type Sep struct {
	mid, last string
	lines     []string
}

// AndSeparator creates a phrase: a, b, c, and d.
// Note: spacing between words is left to printer.Span.
func AndSeparator() *Sep {
	return &Sep{mid: ",", last: ", and"}
}

// OrSeparator creates a phrase: a, b, c, or d.
// Note: spacing between words is left to printer.Span.
func OrSeparator() *Sep {
	return &Sep{mid: ",", last: ", or"}
}

// WriteTo writes all current lines, with appropriate separators, to the passed output.
// FIX? a slightly better interface might be to pass a writer to Sep, write a little bit -- delayed by one -- at a time, and use io.Closer to flush.
func (l *Sep) WriteTo(w io.Writer) (ret int64, err error) {
	if cnt := len(l.lines); cnt > 0 {
		write := func(s string) bool {
			r, e := io.WriteString(w, s)
			ret += int64(r)
			err = e
			return err == nil
		}
		for i := 0; i < cnt; i++ {
			if i != 0 {
				var sep string
				if more := i+1 < cnt; more {
					sep = l.mid
				} else {
					sep = l.last
				}
				if !write(sep) {
					break
				}
			}
			if !write(l.lines[i]) {
				break
			}
		}
	}
	return
}

// Write implements io.Writer, treating every Write as a new line.
func (l *Sep) Write(p []byte) (ret int, err error) {
	var buf bytes.Buffer
	if r, e := buf.Write(p); e != nil {
		err = e
	} else {
		l.lines = append(l.lines, buf.String())
		ret = r
	}
	return
}
