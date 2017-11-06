package postfix

import (
	"fmt"
	"strings"
)

// Function element of an Expression.
type Function interface {
	// Name of the function.
	// Should be a constant string for each function instance.
	Name() string
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
		s[i] = fmt.Sprint(f)
	}
	return strings.Join(s, " ")
}
