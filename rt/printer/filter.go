package printer

import (
	"github.com/ionous/iffy/lang"
	"io"
	"strings"
)

// Bracket filters io.Writer, parenthesizing a stream of writes.
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

// Close implements io.Closer to terminate the parenthesis. It does not close the wrapped stream.
func (l *Bracket) Close() (err error) {
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

func (l *TitleCase) Write(p []byte) (int, error) {
	return io.WriteString(l.Writer, lang.Capitalize(string(p)))
}
