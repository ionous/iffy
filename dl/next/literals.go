package next

import (
	"io"
	"strconv"

	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/qna"
	"github.com/ionous/iffy/rt"
)

// BoolValue specifies a simple true/false value.
type BoolValue struct {
	Bool bool
}

// Compose returns a spec for use by the composer editor.
func (*BoolValue) Compose() composer.Spec {
	return composer.Spec{
		Name:  "bool_value",
		Spec:  "{bool:bool_eval}",
		Group: "literals",
		Desc:  "Bool Value: specify an explicit true or false value.",
	}
}

// GetBool implements BoolEval; providing the dl with a boolean literal.
func (b *BoolValue) GetBool(rt.Runtime) (bool, error) {
	return b.Bool, nil
}

// String uses strconv.FormatBool.
func (b *BoolValue) String() string {
	return strconv.FormatBool(b.Bool)
}

// NumValue specifies a number value.
type NumValue struct {
	Num float64
}

// GetNumber implements NumberEval providing the dl with a number literal.
func (n *NumValue) GetNumber(rt.Runtime) (float64, error) {
	return n.Num, nil
}

// Int converts to native int.
func (n *NumValue) Int() int {
	return int(n.Num)
}

// Float converts to native float.
func (n *NumValue) Float() float64 {
	return n.Num
}

// String returns a nicely formatted float, with no decimal point when possible.
func (n *NumValue) String() string {
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
func (t *Text) WriteText(run rt.Runtime, w io.Writer) error {
	_, e := io.WriteString(w, t.Text)
	return e
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

// ObjectNames specifies multiple object names.
type ObjectNames struct {
	Names []string
}

// GetObjectStream returns a stream of names without bothering to validate if they exist
func (l *ObjectNames) GetObjectStream(run rt.Runtime) (rt.ObjectStream, error) {
	return qna.NewObjectList(l.Names), nil
}
