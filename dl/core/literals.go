package core

import (
	"strconv"

	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

// Bool specifies a simple true/false value.
type Bool struct {
	Bool bool
}

// Number specifies a number value.
type Number struct {
	Num float64
}

// Text specifies a string value.
type Text struct {
	Text string
}

// Lines specifies a potentially multi-line string value.
type Lines struct {
	Lines string
}

// Numbers specifies multiple float values.
type Numbers struct {
	Values []float64
}

// Texts specifies multiple strings.
type Texts struct {
	Values []string
}

// Compose returns a spec for use by the composer editor.
func (*Bool) Compose() composer.Spec {
	return composer.Spec{
		Name:  "bool_value",
		Spec:  "{bool}",
		Group: "literals",
		Desc:  "Bool Value: specify an explicit true or false value.",
	}
}

// GetBool implements BoolEval; providing the dl with a boolean literal.
func (op *Bool) GetBool(rt.Runtime) (ret g.Value, _ error) {
	ret = g.BoolOf(op.Bool)
	return
}

// String uses strconv.FormatBool.
func (op *Bool) String() string {
	return strconv.FormatBool(op.Bool)
}

func (*Number) Compose() composer.Spec {
	return composer.Spec{
		Name:  "num_value",
		Group: "literals",
		Spec:  "{num:number}",
		Desc:  "Number Value: Specify a particular number.",
	}
}

// GetNumber implements NumberEval providing the dl with a number literal.
func (op *Number) GetNumber(rt.Runtime) (ret g.Value, _ error) {
	ret = g.FloatOf(op.Num)
	return
}

// Int converts to native int.
func (op *Number) Int() int {
	return int(op.Num)
}

// Float converts to native float.
func (op *Number) Float() float64 {
	return op.Num
}

// String returns a nicely formatted float, with no decimal point when possible.
func (op *Number) String() string {
	return strconv.FormatFloat(op.Num, 'g', -1, 64)
}

func (*Text) Compose() composer.Spec {
	return composer.Spec{
		Name:  "text_value",
		Spec:  "{text}",
		Group: "literals",
		Desc:  "Text Value: specify a small bit of text.",
		Stub:  true,
	}
}

// GetText implements interface TextEval providing the dl with a text literal.
func (op *Text) GetText(run rt.Runtime) (ret g.Value, _ error) {
	ret = g.StringOf(op.Text)
	return
}

// String returns the text.
func (op *Text) String() string {
	return op.Text
}

func (*Lines) Compose() composer.Spec {
	return composer.Spec{
		Name:  "lines_value",
		Spec:  "{lines|quote}",
		Group: "literals",
		Desc:  "Lines Value: specify one or more lines of text.",
	}
}

// GetLines implements interface LinesEval providing the dl with a lines literal.
func (op *Lines) GetLines(run rt.Runtime) (ret g.Value, _ error) {
	ret = g.StringOf(op.Lines)
	return
}

// String returns the lines.
func (op *Lines) String() string {
	return op.Lines
}

func (*Numbers) Compose() composer.Spec {
	return composer.Spec{
		Group: "literals",
		Desc:  "Number List: Specify a list of multiple numbers.",
	}
}

func (op *Numbers) GetNumList(rt.Runtime) (ret g.Value, _ error) {
	// note: this generates a new slice pointing to the op.Values memory;
	// fix: should this be a copy? or, maybe mark this as read-only
	ret = g.FloatsOf(op.Values)
	return
}

func (*Texts) Compose() composer.Spec {
	return composer.Spec{
		Group: "literals",
		Desc:  "Text List: specifies multiple string values.",
		Spec:  "text {values*text|comma-and}",
	}
}

func (op *Texts) GetTextList(rt.Runtime) (ret g.Value, _ error) {
	ret = g.StringsOf(op.Values)
	return
}
