package chart

import (
	"fmt"
	"github.com/ionous/iffy/template/postfix"
)

// Parse the passed template string into an expression.
// This is package chart's primary method.
func Parse(template string) (ret postfix.Expression, err error) {
	p := MakeTemplateParser()
	if e := parse(&p, template); e != nil {
		if p.err != nil {
			err = p.err
		} else {
			err = e
		}
	} else {
		ret, err = p.GetExpression()
	}
	return
}

// ParseExpression reads a series of simple operand and operator phrases.
func ParseExpression(str string) (ret postfix.Expression, err error) {
	var p SeriesParser
	if e := parse(&p, str); e != nil {
		err = e
	} else {
		ret, err = p.GetExpression()
	}
	return
}

type endpointError struct {
	str  string
	end  int
	last State
}

func (e endpointError) Error() string {
	return fmt.Sprintf("parsing of `%s` ended in %T at %q(%d)",
		e.str, e.last, e.str[e.end-1], e.end)
}

func parse(try State, str string) (err error) {
	if end, last := innerParse(try, str); end == 0 {
		last.NewRune(eof) // FIX: this is odd, i know.
	} else {
		err = endpointError{str, end, last}
	}
	return
}

// if the state machine ends before the string is empty
// return the one-index point of failure; 0 therefore is "ok".
func innerParse(try State, str string) (ret int, last State) {
	for i, r := range str {
		if next := try.NewRune(r); next != nil {
			try = next
		} else {
			ret = i + 1
			break
		}
	}
	last = try
	return
}
