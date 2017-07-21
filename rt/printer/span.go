package printer

import (
	"bytes"
	"unicode"
)

// Span implements io.Writer, treating each new Write as word and adding spaces to separate words as necessary.
type Span struct {
	buf bytes.Buffer
}

// String returns the accumulated words as a string.
func (p *Span) String() string {
	return p.buf.String()
}

// Bytes returns the accumulated words as an array of bytes.
func (p *Span) Bytes() []byte {
	return p.buf.Bytes()
}

// Write implements io.Writer treading writes as words.
// ex. Writing "hello", "there,", "world." becomes "hello there. world."
func (p *Span) Write(s []byte) (ret int, err error) {
	// printed something before?
	if len(s) > 0 {
		if p.buf.Len() > 0 {
			// before writing this new thing, possibly put a space.
			letter := []rune(string(s))[0]
			if !unicode.IsSpace(letter) && !unicode.In(letter, unicode.Po, unicode.Pi, unicode.Pf) {
				n, _ := p.buf.WriteString(" ")
				ret += n
			}
		}
		n, _ := p.buf.Write(s)
		ret += n
	}
	return
}

// https://www.compart.com/en/unicode/category/Pi
// Pc     = _Pc // Pc is the set of Unicode characters in category Pc.
// Pd     = _Pd // Pd is the set of Unicode characters in category Pd.
// Pe     = _Pe // Pe is the set of Unicode characters in category Pe.
// Pf     = _Pf // Pf is the set of Unicode characters in category Pf.
// Pi     = _Pi // Pi is the set of Unicode characters in category Pi.
// Po     = _Po // Po is the set of Unicode characters in category Po.
// Ps     = _Ps // Ps is the set of Unicode characters in category Ps.
//
// connector: _
// dash: -|
// end/close: )〞-- closing brackets, some arabic? right ornate quotes
// final: ’” -- angled right quotes
// initial: ‘“ -- angled left quotes
// other: !"'.; -- straight quotes, and full stop
// start/open: ([{ -- opening brackets, some arabic? left ornate quotes
