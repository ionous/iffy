package print

import (
	"bytes"
	"unicode"

	"github.com/ionous/iffy/rt/writer"
)

// Spanner buffers with spacing, treating each new Write as word and adding spaces to separate words as necessary.
type Spanner struct {
	writer.ChunkOutput
	buf bytes.Buffer // note: we cant aggregate buf or io.WriteString will bypasses implementation of Write() in favor of bytes.Buffer.WriteString()
}

func NewSpanner() *Spanner {
	s := new(Spanner)
	writer.InitChunks(s)
	return s
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

func (p *Spanner) WriteChunk(c writer.Chunk) (ret int, err error) {
	// writing something?
	if b, cnt := c.DecodeRune(); cnt > 0 {
		// and already written something and the thing we are writing is not a space?
		if p.Len() > 0 && !spaceLike(b) {
			n, _ := p.buf.WriteRune(' ')
			ret += n
		}
		n, _ := c.WriteTo(&p.buf)
		ret += n
	}
	return
}

func spaceLike(r rune) bool {
	return unicode.IsSpace(r) || unicode.In(r, unicode.Po, unicode.Pi, unicode.Pf)
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
