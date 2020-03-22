package postfix

import (
	"strings"
)

// Function element of an Expression.
// Examples of function(s) are mathematical and boolean operators, if,else blocks, iffy commands, etc.
type Function interface {
	// Arity the number of required function arguments;
	// Should be a constant number for each function instance.
	Arity() int
	// Precedence for infix priority; only has meaning where Arity is non-zero.
	// Should be a constant number for each function instance.
	Precedence() int
}

// Expression provides a series of post fix functions:
// even primitive values are represented here as zero-arity functions returning the value in question.
type Expression []Function

func (x Expression) String() string {
	s := make([]string, len(x))
	for i, f := range x {
		if str, ok := f.(stringer); !ok {
			s[i] = "???"
		} else {
			s[i] = str.String()
		}
	}
	return strings.Join(s, " ")
}

type stringer interface {
	String() string
}
