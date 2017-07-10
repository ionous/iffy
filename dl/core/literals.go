package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
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
	return NewNumberStream(l.Values), nil
}

func NewNumberStream(list []float64) rt.NumberStream {
	return &NumberIt{list: list}
}

type NumberIt struct {
	list []float64
	idx  int
}

func (it *NumberIt) HasNext() bool {
	return it.idx < len(it.list)
}

func (it *NumberIt) GetNext() (ret float64, err error) {
	if !it.HasNext() {
		err = rt.StreamExceeded
	} else {
		ret = it.list[it.idx]
		it.idx++
	}
	return
}

// Texts specifies multiple strings.
type Texts struct {
	Values []string
}

func (l *Texts) GetTextStream(rt.Runtime) (rt.TextStream, error) {
	return &TextIt{list: l.Values}, nil
}

func NewTextStream(list []string) rt.TextStream {
	return &TextIt{list: list}
}

type TextIt struct {
	list []string
	idx  int
}

func (it *TextIt) HasNext() bool {
	return it.idx < len(it.list)
}

func (it *TextIt) GetNext() (ret string, err error) {
	if !it.HasNext() {
		err = rt.StreamExceeded
	} else {
		ret = it.list[it.idx]
		it.idx++
	}
	return
}

// Objects specifies multiple object names.
type Objects struct {
	Names []string
}

func (l *Objects) GetObjectStream(run rt.Runtime) (rt.ObjectStream, error) {
	return &ObjectIt{run: run, list: l.Names}, nil
}

func NewObjectStream(run rt.Runtime, list []string) rt.ObjectStream {
	return &ObjectIt{run: run, list: list}
}

type ObjectIt struct {
	run  rt.Runtime
	list []string
	idx  int
}

func (it *ObjectIt) HasNext() bool {
	return it.idx < len(it.list)
}

func (it *ObjectIt) GetNext() (ret rt.Object, err error) {
	if !it.HasNext() {
		err = rt.StreamExceeded
	} else {
		ref := it.list[it.idx]
		if obj, ok := it.run.FindObject(ref); !ok {
			err = errutil.New("couldnt find object", ref)
		} else {
			ret = obj
			it.idx++
		}
	}
	return
}
