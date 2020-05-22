package chart

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/template/postfix"
	"github.com/ionous/iffy/template/types"
)

// implements OperandState.
type QuoteParser struct {
	runes Runes
	err   error
}

func (p *QuoteParser) StateName() string {
	return "quotes"
}

// NewRune starts with the leading quote mark; it finishes just after the matching quote mark.
func (p *QuoteParser) NewRune(r rune) (ret State) {
	if isQuote(r) {
		ret = p.scanQuote(r)
	}
	return
}

func (p *QuoteParser) GetOperand() (ret postfix.Function, err error) {
	if r, e := p.GetString(); e != nil {
		err = e
	} else {
		ret = types.Quote(r)
	}
	return
}

// GetString returns the text including its surrounding quote markers.
func (p *QuoteParser) GetString() (ret string, err error) {
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
	return SelfStatement("findMatchingQuote", func(self State, r rune) (ret State) {
		switch {
		case r == q:
			// eats the ending quote
			ret = Terminal

		case r == escape:
			ret = Statement("escape", func(r rune) (ret State) {
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
