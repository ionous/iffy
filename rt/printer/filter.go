package printer

import (
	"github.com/ionous/iffy/lang"
	"io"
	"strings"
)

// Bracket filters io.Writer, parenthesizing a stream of writes. Close adds the closing paren.
func Bracket(w io.Writer) io.WriteCloser {
	return &_Bracket{Writer: w}
}

// Capitalize filters io.Writer, capitalizing the first string.
func Capitalize(w io.Writer) io.Writer {
	return &_Capitalize{Writer: w}
}

// Lowercase filters io.Writer, lowering every string.
func Lowercase(w io.Writer) io.Writer {
	return &_Lowercase{Writer: w}
}

// Spanning filters io.Writer as per Span, writing the final result to the passed buffer.
func Spanning(w io.Writer) io.Writer {
	return &_Spanning{w: w}
}

// TitleCase filters io.Writer, capitalizing every write.
func TitleCase(w io.Writer) io.Writer {
	return &_TitleCase{Writer: w}
}

type _Bracket struct {
	io.Writer
	cnt int
}

type _Capitalize struct {
	io.Writer
	cnt int
}

type _Lowercase struct {
	io.Writer
}

type _Spanning struct {
	Span
	w io.Writer
}

type _TitleCase struct {
	io.Writer
}

func (l *_Bracket) Write(p []byte) (ret int, err error) {
	if l.cnt == 0 {
		io.WriteString(l.Writer, "(")
	}
	l.cnt++
	return l.Writer.Write(p)
}

// Close to terminate the parenthesis.
func (l *_Bracket) Close() (err error) {
	if l.cnt > 0 {
		_, err = io.WriteString(l.Writer, ")")
	}
	return
}

func (l *_Capitalize) Write(p []byte) (ret int, err error) {
	if l.cnt == 0 {
		ret, err = io.WriteString(l.Writer, lang.Capitalize((string(p))))
	} else {
		ret, err = l.Writer.Write(p)
	}
	l.cnt++
	return
}

func (l *_Lowercase) Write(p []byte) (int, error) {
	return io.WriteString(l.Writer, strings.ToLower((string(p))))
}

func (l *_Spanning) Close() (err error) {
	if l.Len() > 0 {
		_, err = l.w.Write(l.Bytes())
	}
	return
}

func (l *_TitleCase) Write(p []byte) (int, error) {
	return io.WriteString(l.Writer, lang.Capitalize(string(p)))
}
