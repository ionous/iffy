package kind

import (
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

// Number returns reflect.TypeOf(float64)
func Number() r.Type { return number }

// Number returns reflect.TypeOf(float64)
func IdentId() r.Type { return identId }

// BoolEval returns reflect.TypeOf(rt.BoolEval)
func BoolEval() r.Type { return boolEval }

// NumberEval returns reflect.TypeOf(rt.NumberEval)
func NumberEval() r.Type { return numEval }

// TextEval returns reflect.TypeOf(rt.TextEval)
func TextEval() r.Type { return textEval }

// ObjectEval returns reflect.TypeOf(rt.ObjectEval)
func ObjectEval() r.Type { return objEval }

// NumListEval returns reflect.TypeOf(rt.NumListEval)
func NumListEval() r.Type { return numListEval }

// TextListEval returns reflect.TypeOf(rt.TextListEval)
func TextListEval() r.Type { return textListEval }

// ObjListEval returns reflect.TypeOf(rt.ObjListEval)
func ObjListEval() r.Type { return objListEval }

// switches dont work well with .Interface().(type) when dst is nil.
var number r.Type = r.TypeOf((*float64)(nil)).Elem()
var identId = r.TypeOf((*ident.Id)(nil)).Elem()
var boolEval = r.TypeOf((*rt.BoolEval)(nil)).Elem()
var numEval = r.TypeOf((*rt.NumberEval)(nil)).Elem()
var textEval = r.TypeOf((*rt.TextEval)(nil)).Elem()
var objEval = r.TypeOf((*rt.ObjectEval)(nil)).Elem()
var numListEval = r.TypeOf((*rt.NumListEval)(nil)).Elem()
var textListEval = r.TypeOf((*rt.TextListEval)(nil)).Elem()
var objListEval = r.TypeOf((*rt.ObjListEval)(nil)).Elem()

// Number
// BoolEval
// NumberEval
// TextEval
// ObjectEval
// NumListEval
// TextListEval
// ObjListEval
