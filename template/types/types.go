package types

import (
	"strconv"
	"strings"
)

// Argument, a marker for zero-artiy functions.
type Argument interface {
	arg()
}

// Quote, a zero-arity function returning a string.
type Quote string

// Number, a zero-arity function returning a number.
type Number float64

// Bool, a zero-arity function returning a true/false value.
type Bool bool

// Reference, a zero-arity function containing the path to an object property.
type Reference []string

// Command, a dynamic-arity function referring to some externally defined command.
type Command struct {
	CommandName  string
	CommandArity int
}

// Builtin, a dynamic-arity function referring to one of a small set predefined commands.
// See also: BuiltinType.
type Builtin struct {
	Type           BuiltinType
	ParameterCount int
}

type BuiltinType int

//go:generate stringer -type=BuiltinType
const (
	// Joins several expressions into a single result.
	// For example, if multiple strings live side by side
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
