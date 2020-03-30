package printer

import (
	"io"
	"strings"

	"github.com/ionous/iffy/lang"
)

// Bracket filters io.Writer, parenthesizing a stream of writes. Close adds the closing paren.
func Bracket(w io.Writer) io.WriteCloser {
	return &_Bracket{w: w}
}

// Capitalize filters io.Writer, capitalizing the first string.
func Capitalize(w io.Writer) io.Writer {
	return &_Capitalize{w: w}
}

// Lowercase filters io.Writer, lowering every string.
func Lowercase(w io.Writer) io.Writer {
	return &_Lowercase{w: w}
}

// Slash filters io.Writer, separating writes with a slash.
func Slash(w io.Writer) io.Writer {
	return &_Slash{w: w}
}

// Spacing filters io.Writer as per Span, writing the final result to the passed buffer.
func Spacing(w io.Writer) io.Writer {
	return &_Spacing{w: w}
}

// TitleCase filters io.Writer, capitalizing every write.
func TitleCase(w io.Writer) io.Writer {
	return &_TitleCase{w: w}
}

type _Bracket struct {
	w   io.Writer
	cnt int
}

type _Capitalize struct {
	w   io.Writer
	cnt int
}

type _Lowercase struct {
	w io.Writer
}

type _Slash struct {
	w   io.Writer
	cnt int
}

type _Spacing struct {
	Spanner
	w io.Writer
}

type _TitleCase struct {
	w io.Writer
}

func (l *_Bracket) Write(p []byte) (ret int, err error) {
	if l.cnt == 0 {
		n, e := io.WriteString(l.w, "(")
		l.cnt += n
		err = e
	}
	if err == nil {
		ret, err = l.w.Write(p)
		l.cnt += ret
	}
	return

}

// Close to terminate the parenthesis.
func (l *_Bracket) Close() (err error) {
	if l.cnt > 0 {
		_, err = io.WriteString(l.w, ")")
	}
	return
}

func (l *_Capitalize) Write(p []byte) (ret int, err error) {
	if l.cnt == 0 {
		ret, err = io.WriteString(l.w, lang.Capitalize((string(p))))
	} else {
		ret, err = l.w.Write(p)
	}
	l.cnt++
	return
}

func (l *_Lowercase) Write(p []byte) (int, error) {
	return io.WriteString(l.w, strings.ToLower((string(p))))
}

func (l *_Slash) Write(p []byte) (ret int, err error) {
	if l.cnt != 0 {
		n, _ := io.WriteString(l.w, " /")
		l.cnt += n
	}
	if err == nil {
		ret, err = l.w.Write(p)
		l.cnt += ret
	}
	return
}

func (l *_Spacing) Close() (err error) {
	if l.Len() > 0 {
		_, err = l.w.Write(l.Bytes())
	}
	return
}

func (l *_TitleCase) Write(p []byte) (int, error) {
	return io.WriteString(l.w, lang.Capitalize(string(p)))
}
