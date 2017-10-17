package chart

import (
	"github.com/ionous/errutil"
	"strconv"
)

type Mode int

const (
	Int10 Mode = iota
	Int16
	Float64
	Exponent
)

type numParser struct {
	runes     []rune
	mode      Mode
	suspended bool
	err       error
}

func (p numParser) GetValue() (ret float64, err error) {
	if p.err != nil {
		err = p.err
	} else if p.suspended {
		err = errutil.New("incomplete number")
	} else {
		s := string(p.runes)
		switch p.mode {
		case Int10:
			if i, e := strconv.ParseInt(s, 10, 64); e != nil {
				panic(e)
			} else {
				ret = float64(i)
			}
		case Int16:
			if i, e := strconv.ParseInt(s[2:], 16, 64); e != nil {
				panic(e)
			} else {
				ret = float64(i)
			}
		case Float64, Exponent:
			if f, e := strconv.ParseFloat(s, 64); e != nil {
				panic(e)
			} else {
				ret = f
			}
		default:
			err = errutil.New("unknown number", s, p.mode)
		}
	}
	return
}

// initial state of digit parsing, moves to or delegates to firstDigit
// note: iffy doesn't support leading with just a "."
func (p *numParser) NewRune(r rune) (ret State) {
	// in go, the leading +/- are unary operators;
	// and not pieces of the number
	if r == '-' || r == '+' {
		ret = p.maybe(r, p.firstDigit)
	} else {
		ret = p.firstDigit(r)
	}
	return
}

// helper to improve readability: accepts a rune as part of a number.
func (p *numParser) yes(r rune, s Statement) State {
	p.update(r, true)
	return s
}

// helper to improve readability: provisionally accepts a rune as part of a number.
func (p *numParser) maybe(r rune, s Statement) State {
	p.update(r, false)
	return s
}

func (p *numParser) update(r rune, active bool) {
	p.runes = append(p.runes, r)
	p.suspended = !active
}

func (p *numParser) firstDigit(r rune) (ret State) {
	if r >= '1' && r <= '9' {
		// https://golang.org/ref/spec#float_lit
		ret = p.yes(r, p.decimalDigit)
	} else if r == '0' {
		// 0 can standalone; but, it might be followed by a hex qualifier.
		ret = p.yes(r, func(r rune) (ret State) {
			// https://golang.org/ref/spec#hex_literal
			if r == 'x' || r == 'X' {
				p.mode = Int16
				ret = p.maybe(r, p.hexDigits)
			} else {
				ret = p.decimalDigit(r)
			}
			return
		})
	}
	return
}

// a string of numbers, possibly followed by a decimal or exponent separator.
// note: golang numbers can end in a pure ".", iffy chooses not to allow that.
func (p *numParser) decimalDigit(r rune) (ret State) {
	if isNumber(r) {
		ret = p.yes(r, p.decimalDigit)
	} else if r == '.' {
		p.mode = Float64
		ret = p.maybe(r, func(r rune) (ret State) {
			if isNumber(r) {
				ret = p.yes(r, p.decimalDigit)
			} else {
				ret = p.tryExponent(r)
			}
			return
		})
	} else {
		ret = p.tryExponent(r)
	}
	return
}

// https://golang.org/ref/spec#exponent
// exponent  = ( "e" | "E" ) [ "+" | "-" ] decimals
func (p *numParser) tryExponent(r rune) (ret State) {
	if r == 'e' || r == 'E' {
		p.mode = Exponent
		ret = p.maybe(r, func(r rune) (ret State) {
			if r == '+' || r == '-' {
				ret = p.maybe(r, p.decimals)
			} else if isNumber(r) {
				ret = p.yes(r, p.decimals)
			}
			return
		})
	}
	return
}

// a chain of decimal digits 0-9
func (p *numParser) decimals(r rune) (ret State) {
	if isNumber(r) {
		ret = p.yes(r, p.decimals)
	}
	return
}

// a chain of hex digits 0-9, a-f
func (p *numParser) hexDigits(r rune) (ret State) {
	if isHex(r) {
		ret = p.yes(r, p.hexDigits)
	}
	return
}
