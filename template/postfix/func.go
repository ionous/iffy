package postfix

import (
	"strings"
)

// Function element of an Expression.
type Function interface {
	// Arity the number of required function arguments;
	// Should be a constant number for each function instance.
	Arity() int
	// Precedence for infix priority; only has meaning where Arity is non-zero.
	// Should be a constant number for each function instance.
	Precedence() int
}

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
