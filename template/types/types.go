package types

import (
	"strconv"
	"strings"
)

// Argument, a marker for zero-artiy functions.
type Argument interface {
	arg()
}

// Quote text as a postfix function.
type Quote string

// Number as a postfix function.
type Number float64

// Bool as a postfix function.
type Bool bool

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

// Builtin
type Builtin struct {
	Type           BuiltinType
	ParameterCount int
}

type BuiltinType int

//go:generate stringer -type=BuiltinType
const (
	// Joins several expressions into a single result.
	// Usually, each individual expression reduces to a string,
	// and span then concatenates those strings.
	Span BuiltinType = iota
	// sequences holds a set of potentially repeating text.
	Stopping
	Shuffle
	Cycle
	// arity varies based on whether there's an else statement
	IfStatement
	UnlessStatement
)

func (Quote) Arity() int      { return 0 }
func (Quote) Precedence() int { return 0 }
func (k Quote) Value() string { return string(k) }

func (Number) Arity() int       { return 0 }
func (Number) Precedence() int  { return 0 }
func (k Number) Value() float64 { return float64(k) }

func (Bool) Arity() int      { return 0 }
func (Bool) Precedence() int { return 0 }
func (k Bool) Value() bool   { return bool(k) }

func (Reference) Arity() int        { return 0 }
func (Reference) Precedence() int   { return 0 }
func (r Reference) Value() []string { return []string(r) }

func (k Command) Arity() int      { return k.CommandArity }
func (k Command) Precedence() int { return 8 }

func (k Builtin) Arity() int      { return k.ParameterCount }
func (k Builtin) Precedence() int { return 8 }

func (k Quote) String() string {
	return strconv.Quote(string(k))
}
func (k Number) String() string {
	return strconv.FormatFloat(float64(k), 'g', -1, 64)
}
func (k Bool) String() string {
	return strconv.FormatBool(bool(k))
}
func (k Reference) String() string {
	return strings.Join([]string(k), ".")
}
func (k Command) String() string {
	w := strings.ToUpper(k.CommandName)
	a := strconv.Itoa(k.CommandArity)
	return w + "/" + a
}
func (k Builtin) String() string {
	w := k.Type.String()
	a := strconv.Itoa(k.ParameterCount)
	return w + "/" + a
}
