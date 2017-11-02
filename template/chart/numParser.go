package chart

import (
	"github.com/ionous/errutil"
	"strconv"
)

type FloatMode int

//go:generate stringer -type=FloatMode
const (
	Pending FloatMode = iota
	Int10
	Int16
	Float64
)

type NumParser struct {
	runes  Runes
	mode   FloatMode
	negate negate
}

type negate bool

func (n negate) mul(f float64) float64 {
	if n {
		f *= -1.0
	}
	return f
}

func (p NumParser) GetValue() (ret float64, err error) {
	s := p.runes.String()
	switch p.mode {
	case Int10:
		if i, e := strconv.ParseInt(s, 10, 64); e != nil {
			panic(e)
		} else {
			ret = p.negate.mul(float64(i))
		}
	case Int16: // chops out the 0x qualifier
		if i, e := strconv.ParseInt(s[2:], 16, 64); e != nil {
			panic(e)
		} else {
			ret = float64(i)
		}
	case Float64:
		if f, e := strconv.ParseFloat(s, 64); e != nil {
			panic(e)
		} else {
			ret = p.negate.mul(f)
		}
	default:
		err = errutil.Fmt("unknown number: '%v' is %v.", s, p.mode)
	}
	return
}

// initial state of digit parsing.
// note: iffy doesn't support leading with just a "."
func (p *NumParser) NewRune(r rune) (ret State) {
	switch {
	// in golang, leading +/- are unary operators;
	// in iffy, they are considered optional parts decimal numbers.
	// note: strconv's base 10 parser doesnt handle leading signs.
	// we therefore leave them out of our result, and just flag the negative ones.
	case r == '-':
		p.negate = true
		fallthrough
	case r == '+':
		ret = Statement(func(r rune) (ret State) {
			if isNumber(r) {
				p.mode = Int10
				ret = p.runes.Accept(r, Statement(p.leadingDigit))
			}
			return
		})
	case r == '0':
		// 0 can standalone; but, it might be followed by a hex qualifier.
		p.mode = Int10
		ret = p.runes.Accept(r, Statement(func(r rune) (ret State) {
			// https://golang.org/ref/spec#hex_literal
			switch {
			case r == 'x' || r == 'X':
				p.mode = Pending
				ret = p.runes.Accept(r, Statement(func(r rune) (ret State) {
					if isHex(r) {
						p.mode = Int16
						ret = p.runes.Accept(r, Statement(p.hexDigits))
					}
					return
				}))
			default:
				// delegate to number and dot checking...
				// in a statechart, it would be a super-state, and
				// x (above) would jump to a sibling of that super-state.
				ret = p.leadingDigit(r)
			}
			return
		}))
	case isNumber(r):
		// https://golang.org/ref/spec#float_lit
		p.mode = Int10
		ret = p.runes.Accept(r, Statement(p.leadingDigit))
	}
	return
}

// a string of numbers, possibly followed by a decimal or exponent separator.
// note: golang numbers can end in a pure ".", iffy chooses not to allow that.
func (p *NumParser) leadingDigit(r rune) (ret State) {
	switch {
	case isNumber(r):
		ret = p.runes.Accept(r, Statement(p.leadingDigit))
	case r == '.':
		p.mode = Pending
		ret = p.runes.Accept(r, Statement(func(r rune) (ret State) {
			if isNumber(r) {
				p.mode = Float64
				ret = p.runes.Accept(r, Statement(p.leadingDigit))
			} else {
				ret = p.tryExponent(r) // delegate to exponent checking,,,
			}
			return
		}))
	default:
		ret = p.tryExponent(r) // delegate to exponent checking,,,
	}
	return
}

// https://golang.org/ref/spec#exponent
// exponent  = ( "e" | "E" ) [ "+" | "-" ] decimals
func (p *NumParser) tryExponent(r rune) (ret State) {
	switch {
	case r == 'e' || r == 'E':
		p.mode = Pending
		ret = p.runes.Accept(r, Statement(func(r rune) (ret State) {
			switch {
			case isNumber(r):
				p.mode = Float64
				ret = p.runes.Accept(r, Statement(p.decimals))
			case r == '+' || r == '-':
				ret = p.runes.Accept(r, Statement(func(r rune) (ret State) {
					if isNumber(r) {
						p.mode = Float64
						ret = p.runes.Accept(r, Statement(p.decimals))
					}
					return
				}))
			}
			return
		}))
	}
	return
}

// a chain of decimal digits 0-9
func (p *NumParser) decimals(r rune) (ret State) {
	if isNumber(r) {
		ret = p.runes.Accept(r, Statement(p.decimals))
	}
	return
}

// a chain of hex digits 0-9, a-f
func (p *NumParser) hexDigits(r rune) (ret State) {
	if isHex(r) {
		ret = p.runes.Accept(r, Statement(p.hexDigits))
	}
	return
}
