package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
	"strconv"
)

// Bool specifies a simple true/false value.
type Bool struct {
	Bool bool
}

// GetBool implements BoolEval; providing the dl with a boolean literal.
func (b *Bool) GetBool(rt.Runtime) (bool, error) {
	return b.Bool, nil
}

// String uses strconv.FormatBool.
func (b *Bool) String() string {
	return strconv.FormatBool(b.Bool)
}

// Num specifies a number value.
type Num struct {
	Num float64
}

// GetNumber implements NumberEval providing the dl with a number literal.
func (n *Num) GetNumber(rt.Runtime) (float64, error) {
	return n.Num, nil
}

// Int converts to native int.
func (n *Num) Int() int {
	return int(n.Num)
}

// Float converts to native float.
func (n *Num) Float() float64 {
	return n.Num
}

// String returns a nicely formatted float, with no decimal point when possible.
func (n *Num) String() string {
	return strconv.FormatFloat(n.Num, 'g', -1, 64)
}

// Text specifies a string value.
type Text struct {
	Text string
}

// GetText implements interface TextEval providing the dl with a text literal.
func (t *Text) GetText(rt.Runtime) (string, error) {
	return t.Text, nil
}

// String returns the text.
func (t *Text) String() string {
	return t.Text
}

// Object searches for objects in the world by name.
type Object struct {
	Name string
}

// GetObject searches through the scope for an object matching Name
func (op *Object) GetObject(run rt.Runtime) (ret rt.Object, err error) {
	if obj, ok := run.FindObject(op.Name); !ok {
		err = errutil.New("Object.GetObject, couldnt find", op.Name)
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
	return stream.NewNumberStream(l.Values), nil
}

// Texts specifies multiple strings.
type Texts struct {
	Values []string
}

func (l *Texts) GetTextStream(rt.Runtime) (rt.TextStream, error) {
	return stream.NewTextStream(l.Values), nil
}

// Objects specifies multiple object names.
type Objects struct {
	Names []string
}

func (l *Objects) GetObjectStream(run rt.Runtime) (rt.ObjectStream, error) {
	return stream.NewNameStream(run, l.Names), nil
}
