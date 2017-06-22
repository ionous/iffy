package core

import (
	"bytes"
	"unicode"
)

// SpanPrinter implements io.Writer, treating each new Write as word and adding spaces to separate words as necessary.
type SpanPrinter struct {
	buf bytes.Buffer
}

// String returns the accumulated words as a string.
func (p *SpanPrinter) String() string {
	return p.buf.String()
}

// Bytes returns the accumulated words as an array of bytes.
func (p *SpanPrinter) Bytes() []byte {
	return p.buf.Bytes()
}

// Write implements io.Writer treading writes as words.
// ex. Writing "hello", "there,", "world." becomes "hello there. world."
func (p *SpanPrinter) Write(s []byte) (ret int, err error) {
	// printed something before?
	if len(s) > 0 {
		if p.buf.Len() > 0 {
			// before writing this new thing, possibly put a space.
			letter := []rune(string(s))[0]
			if !unicode.IsPunct(letter) && !unicode.IsSpace(letter) {
				n, _ := p.buf.WriteString(" ")
				ret += n
			}
		}
		if err == nil {
			n, _ := p.buf.Write(s)
			ret += n
		}
	}
	return
}
