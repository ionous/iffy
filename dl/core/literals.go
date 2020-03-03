package core

import (
	"strconv"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
)

// BoolValue specifies a simple true/false value.
type BoolValue struct {
	Bool bool
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

// TextValue specifies a string value.
type TextValue struct {
	Text string `if:"spec:{text:lines|quote}"`
}

// GetText implements interface TextEval providing the dl with a text literal.
func (t *TextValue) GetText(rt.Runtime) (string, error) {
	return t.Text, nil
}

// String returns the text.
func (t *TextValue) String() string {
	return t.Text
}

// TopObject asks for the top most object in scope. It fails if no scope has been established.
type TopObject struct{}

// GetObject searches through the scope for an object matching Name
func (TopObject) GetObject(run rt.Runtime) (ret rt.Object, err error) {
	if obj, ok := run.TopObject(); !ok {
		err = errutil.New("no top object")
	} else {
		ret = obj
	}
	return
}

// ObjectName searches for objects in the world by name.
type ObjectName struct {
	Name string
}

// GetObject searches through the scope for an object matching Name
func (op *ObjectName) GetObject(run rt.Runtime) (ret rt.Object, err error) {
	if obj, ok := run.GetObject(op.Name); !ok {
		err = errutil.New("couldnt find object", op.Name)
	} else {
		ret = obj
	}
	return
}

// Numbers specifies multiple float values.
type Numbers struct {
	Values []float64
}

func (l *Numbers) GetNumberStream(rt.Runtime) (rt.NumberStream, error) {
	return stream.NewNumberStream(stream.FromList(l.Values)), nil
}

// Texts specifies multiple strings.
type Texts struct {
	Values []string
}

func (l *Texts) GetTextStream(rt.Runtime) (rt.TextStream, error) {
	return stream.NewTextStream(stream.FromList(l.Values)), nil
}

// ObjectNames specifies multiple object names.
type ObjectNames struct {
	Names []string
}

func (l *ObjectNames) GetObjectStream(run rt.Runtime) (rt.ObjectStream, error) {
	return stream.NewNameStream(run, l.Names), nil
}
