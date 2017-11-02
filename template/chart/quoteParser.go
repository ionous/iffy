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
		ret = p.runes.Accept(r, p.scanQuote(r))
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
			ret = p.runes.Accept(r, terminal) // done, eat the trailing quote.

		case r == escape:
			ret = Statement(func(r rune) (ret State) {
				if ok := escapes[r]; !ok {
					p.err = errutil.Fmt("unknown escape sequence %q", r)
				} else {
					p.runes.list = append(p.runes.list, escape, r)
					ret = self // loop
				}
				return
			})

		case r != p.runes.list[0]:
			ret = p.runes.Accept(r, self) // loop...
		}
		return
	})
	return
}

var escapes = map[rune]bool{
	'a':  true, // '\a',
	'b':  true, // '\b',
	'f':  true, // '\f',
	'n':  true, // '\n',
	'r':  true, // '\r',
	't':  true, // '\t',
	'v':  true, // '\v',
	'\\': true, //'\\',
	'\'': true, //'\'',
	'"':  true, // '"',
}
