package core

import (
	"github.com/ionous/iffy/rt"
	"strconv"
)

// Bool specifies a simple true/false value.
type Bool struct {
	Value bool
}

// GetBool implements BoolEval; providing the dl with a boolean literal.
func (b *Bool) GetBool(rt.Runtime) (bool, error) {
	return b.Value, nil
}

// String uses strconv.FormatBool.
func (b *Bool) String() string {
	return strconv.FormatBool(b.Value)
}

// Num specifies a number value.
type Num struct {
	Value float64
}

// GetNumber implements NumberEval providing the dl with a number literal.
func (n *Num) GetNumber(rt.Runtime) (float64, error) {
	return n.Value, nil
}

// Int converts to native int.
func (n *Num) Int() int {
	return int(n.Value)
}

// Float converts to native float.
func (n *Num) Float() float64 {
	return n.Value
}

// String returns a nicely formatted float, with no decimal point when possible.
func (n *Num) String() string {
	return strconv.FormatFloat(n.Value, 'g', -1, 64)
}

// Text specifies a string value.
type Text struct {
	Value string
}

// GetText implements interface TextEval providing the dl with a text literal.
func (t *Text) GetText(rt.Runtime) (string, error) {
	return t.Value, nil
}

// String returns the text.
func (t *Text) String() string {
	return t.Value
}
