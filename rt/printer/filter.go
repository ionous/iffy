package printer

import (
	"github.com/ionous/iffy/lang"
	"io"
	"strings"
)

// Bracket filters io.Writer, parenthesizing a stream of writes. Flush adds the closing paren.
type Bracket struct {
	io.Writer
	cnt int
}

// Capitalize filters io.Writer, capitalizing the first string.
type Capitalize struct {
	io.Writer
	cnt int
}

// Lowercase filters io.Writer, lowering every string.
type Lowercase struct {
	io.Writer
}

// Spanner filters io.Writer, accumulating writes via Span. Flush writes the Span to the writer as a single unit.
type Spanner struct {
	Writer io.Writer
	Span
}

// TitleCase filters io.Writer, capitalizing every write.
type TitleCase struct {
	io.Writer
}

func (l *Bracket) Write(p []byte) (ret int, err error) {
	if l.cnt == 0 {
		io.WriteString(l.Writer, "(")
	}
	l.cnt++
	return l.Writer.Write(p)
}

// Flush to terminate the parenthesis.
func (l *Bracket) Flush() (err error) {
	if l.cnt > 0 {
		_, e := io.WriteString(l.Writer, ")")
		err = e
	}
	return
}

func (l *Capitalize) Write(p []byte) (ret int, err error) {
	if l.cnt == 0 {
		ret, err = io.WriteString(l.Writer, lang.Capitalize((string(p))))
	} else {
		ret, err = l.Writer.Write(p)
	}
	l.cnt++
	return
}

func (l *Lowercase) Write(p []byte) (int, error) {
	return io.WriteString(l.Writer, strings.ToLower((string(p))))
}

// Flush to terminate write the accumulated text.
func (l *Spanner) Flush() (err error) {
	if b := l.Bytes(); len(b) > 0 {
		if _, e := l.Writer.Write(b); e != nil {
			err = e
		} else {
			l.Span = Span{}
		}
	}
	return
}

func (l *TitleCase) Write(p []byte) (int, error) {
	return io.WriteString(l.Writer, lang.Capitalize(string(p)))
}
