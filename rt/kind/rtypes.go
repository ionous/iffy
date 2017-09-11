package kind

import (
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

// Bool returns reflect.TypeOf(bool)
func Bool() r.Type { return boolType }

// Int returns reflect.TypeOf(int)
func Int() r.Type { return intType }

// Number returns reflect.TypeOf(float64)
func Number() r.Type { return floatType }

// String returns reflect.TypeOf(string)
func String() r.Type { return stringType }

// IdentId returns reflect.TypeOf(iden.Id)
func IdentId() r.Type { return identId }

// Object returns reflect.TypeOf(rt.Object)
func Object() r.Type { return object }

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
var floatType r.Type = r.TypeOf((*float64)(nil)).Elem()
var stringType r.Type = r.TypeOf((*string)(nil)).Elem()
var intType r.Type = r.TypeOf((*int)(nil)).Elem()
var boolType r.Type = r.TypeOf((*bool)(nil)).Elem()
var object = r.TypeOf((*rt.Object)(nil)).Elem()
var identId = r.TypeOf((*ident.Id)(nil)).Elem()
var boolEval = r.TypeOf((*rt.BoolEval)(nil)).Elem()
var numEval = r.TypeOf((*rt.NumberEval)(nil)).Elem()
var textEval = r.TypeOf((*rt.TextEval)(nil)).Elem()
var objEval = r.TypeOf((*rt.ObjectEval)(nil)).Elem()
var numListEval = r.TypeOf((*rt.NumListEval)(nil)).Elem()
var textListEval = r.TypeOf((*rt.TextListEval)(nil)).Elem()
var objListEval = r.TypeOf((*rt.ObjListEval)(nil)).Elem()

// Bool
// Number
// ident.Id
// BoolEval
// NumberEval
// TextEval
// ObjectEval
// NumListEval
// TextListEval
// ObjListEval
