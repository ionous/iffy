package chart

import (
	"github.com/ionous/errutil"
)

type QuoteParser struct {
	runes Runes
	err   error
}

// NewRune starts with the leading quote mark; it finishes just after the matching quote mark.
func (p *QuoteParser) NewRune(r rune) (ret State) {
	if isQuote(r) {
		ret = p.scanQuote(r)
	}
	return
}

// GetString returns the text including its surrounding quote markers.
func (p QuoteParser) GetString() (ret string, err error) {
	if p.err != nil {
		err = p.err
	} else {
		ret = p.runes.String()
	}
	return
}

// scans until the matching quote marker is found
func (p *QuoteParser) scanQuote(q rune) (ret State) {
	const escape = '\\'
	return SelfStatement(func(self SelfStatement, r rune) (ret State) {
		switch {
		case r == q:
			// for the very next rune returns nil ( unhandled )
			// this rune, the ending quote: it eats.
			ret = Statement(func(rune) State { return nil })

		case r == escape:
			ret = Statement(func(r rune) (ret State) {
				if x, ok := escapes[r]; !ok {
					p.err = errutil.Fmt("unknown escape sequence %q", r)
				} else {
					ret = p.runes.Accept(x, self)
				}
				return
			})

		case r != eof:
			ret = p.runes.Accept(r, self) // loop...
		}
		return
	})
	return
}

var escapes = map[rune]rune{
	'a':  '\a',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'v':  '\v',
	'\\': '\\',
	'\'': '\'',
	'"':  '"',
}
