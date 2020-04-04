package print

import (
	"bytes"
	"io"
	"unicode"
	"unicode/utf8"
)

// Spanner buffers with spacing, treating each new Write as word and adding spaces to separate words as necessary.
type Spanner struct {
	buf bytes.Buffer // note: we cant aggregate buf or io.WriteString will bypasses implementation of Write() in favor of bytes.Buffer.WriteString()
}

func (p *Spanner) Len() int {
	return p.buf.Len()
}
func (p *Spanner) Bytes() []byte {
	return p.buf.Bytes()
}

func (p *Spanner) String() string {
	return p.buf.String()
}

func (p *Spanner) Write(b []byte) (ret int, err error) {
	if len(b) > 0 {
		// printed something before? check for spacing.
		if p.buf.Len() > 0 {
			letter, cnt := utf8.DecodeRune(b)
			if cnt > 0 {
				if !unicode.IsSpace(letter) && !unicode.In(letter, unicode.Po, unicode.Pi, unicode.Pf) {
					n, _ := io.WriteString(&p.buf, " ")
					ret += n
				}
			}
		}
		n, _ := p.buf.Write(b)
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
