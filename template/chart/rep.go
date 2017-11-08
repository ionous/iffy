package chart

import (
	"fmt"
	"strings"
)

// Argument, a marker for zero-artiy functions.
type Argument interface {
	arg()
}

// Quote text as a function.
type Quote string

// Number as a function.
type Number float64

// Reference refers to an object property by name.
// note: there are potentially two ways of representing references:
// a function with "dynamic" arity, each field a parameter;
// a function with 0 arity, and the object is an array of fields.
// choosing the array method for now.
type Reference []string

type Command struct {
	CommandName  string
	CommandArity int
}

func (Quote) Arity() int      { return 0 }
func (Quote) Precedence() int { return 0 }
func (q Quote) Value() string { return string(q) }

func (Number) Arity() int       { return 0 }
func (Number) Precedence() int  { return 0 }
func (n Number) Value() float64 { return float64(n) }

func (Reference) Arity() int        { return 0 }
func (Reference) Precedence() int   { return 0 }
func (r Reference) Value() []string { return []string(r) }

func (c Command) Arity() int      { return c.CommandArity }
func (c Command) Precedence() int { return 8 }

func (q Quote) String() string {
	return fmt.Sprintf("`%s`", string(q))
}
func (n Number) String() string {
	return fmt.Sprintf("%g", float64(n))
}
func (r Reference) String() string {
	return strings.Join([]string(r), ".")
}
func (c Command) String() string {
	return fmt.Sprintf("%s/%d", strings.ToUpper(c.CommandName), c.CommandArity)
}
