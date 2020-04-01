package next

import (
	"strconv"

	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/qna"
	"github.com/ionous/iffy/rt"
)

// Bool specifies a simple true/false value.
type Bool struct {
	Bool bool
}

// Compose returns a spec for use by the composer editor.
func (*Bool) Compose() composer.Spec {
	return composer.Spec{
		Name:  "bool_value",
		Spec:  "{bool:bool_eval}",
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

// Number specifies a number value.
type Number struct {
	Num float64
}

func (*Number) Compose() composer.Spec {
	return composer.Spec{
		Name:  "num_value",
		Group: "literals",
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

// Text specifies a string value.
type Text struct {
	Text string
}

func (*Text) Compose() composer.Spec {
	return composer.Spec{
		Name:  "text_value",
		Spec:  "{text:lines|quote}",
		Group: "literals",
		Desc:  "Text Value: specify one or more lines of text.",
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

// ObjectName searches for objects in the world by name.
type ObjectName struct {
	Name string
}

// GetObject returns the object name if it exists
func (op *ObjectName) GetObject(run rt.Runtime) (ret string, err error) {
	if ok, e := op.GetBool(run); e != nil {
		err = e
	} else if ok {
		ret = op.Name
	}
	return
}

// GetBool returns true if the object exists
func (op *ObjectName) GetBool(run rt.Runtime) (ret bool, err error) {
	e := run.GetObject(op.Name, object.Exists, &ret)
	return ret, e
}

// Numbers specifies multiple float values.
type Numbers struct {
	Values []float64
}

func (l *Numbers) GetNumberStream(rt.Runtime) (rt.NumberStream, error) {
	return qna.NewNumberList(l.Values), nil
}

// Texts specifies multiple strings.
type Texts struct {
	Values []string
}

func (l *Texts) GetTextStream(rt.Runtime) (rt.TextStream, error) {
	return qna.NewTextList(l.Values), nil
}
