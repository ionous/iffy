package chart

import (
	"github.com/ionous/errutil"
)

type quoteParser struct {
	runes []rune
	err   error
}

// GetString returns the text including its surrounding quote markers.
func (p quoteParser) GetString() (ret string, err error) {
	if p.err != nil {
		err = p.err
	} else {
		ret = string(p.runes)
	}
	return
}

// assumes r is the leading quote mark, finish just after the matching quote mark.
func (p *quoteParser) NewRune(r rune) (ret State) {
	if isQuote(r) {
		p.runes = append(p.runes, r)
		ret = p.scanQuote(r)
	}
	return
}

// scans until the matching quote marker is found
func (p *quoteParser) scanQuote(q rune) (ret State) {
	const escape = '\\'
	return SelfStatement(func(self SelfStatement, r rune) (ret State) {
		switch {
		case r == q:
			p.runes = append(p.runes, r)
			ret = terminal // done, eat the trailing quote.

		case r == escape:
			ret = Statement(func(r rune) (ret State) {
				if ok := escapes[r]; !ok {
					p.err = errutil.Fmt("unknown escape sequence %q", r)
				} else {
					p.runes = append(p.runes, escape, r)
					ret = self // loop...
				}
				return
			})

		case r != p.runes[0]:
			p.runes = append(p.runes, r)
			ret = self // loop...
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
