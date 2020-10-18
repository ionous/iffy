package core

import (
	"strconv"

	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
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
		Spec:  "{bool|quote}",
		Group: "literals",
		Desc:  "Bool Value: specify an explicit true or false value.",
	}
}

// GetBool implements BoolEval; providing the dl with a boolean literal.
func (b *Bool) GetBool(rt.Runtime) (bool, error) {
	return b.Bool, nil
}

// String uses strconv.FormatBool.
func (b *Bool) String() string {
	return strconv.FormatBool(b.Bool)
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
func (n *Number) GetNumber(rt.Runtime) (float64, error) {
	return n.Num, nil
}

// Int converts to native int.
func (n *Number) Int() int {
	return int(n.Num)
}

// Float converts to native float.
func (n *Number) Float() float64 {
	return n.Num
}

// String returns a nicely formatted float, with no decimal point when possible.
func (n *Number) String() string {
	return strconv.FormatFloat(n.Num, 'g', -1, 64)
}

func (*Text) Compose() composer.Spec {
	return composer.Spec{
		Name:  "text_value",
		Spec:  "{text}",
		Group: "literals",
		Desc:  "Text Value: specify a small bit of text.",
	}
}

// GetText implements interface TextEval providing the dl with a text literal.
func (t *Text) GetText(run rt.Runtime) (ret string, err error) {
	ret = t.Text
	return
}

// String returns the text.
func (t *Text) String() string {
	return t.Text
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
func (t *Lines) GetLines(run rt.Runtime) (ret string, err error) {
	ret = t.Lines
	return
}

// String returns the lines.
func (t *Lines) String() string {
	return t.Lines
}

func (*Numbers) Compose() composer.Spec {
	return composer.Spec{
		Name:  "numbers",
		Group: "literals",
		Desc:  "Number List: Specify a list of multiple numbers.",
	}
}

func (l *Numbers) GetNumList(rt.Runtime) (ret []float64, _ error) {
	ret = l.Values
	return
}

func (*Texts) Compose() composer.Spec {
	return composer.Spec{
		Name:  "texts",
		Group: "literals",
		Desc:  "Text List: specifies multiple string values.",
	}
}

func (l *Texts) GetTextList(rt.Runtime) (ret []string, _ error) {
	ret = l.Values
	return
}
