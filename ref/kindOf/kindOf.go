package kindOf

import (
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

// Bool kind.
func Bool(rtype r.Type) bool {
	return rtype.Kind() == r.Bool
}

// Int returns true if reflect.Value.Int() would succeed.
func Int(rtype r.Type) (ret bool) {
	switch rtype.Kind() {
	case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
		ret = true
	}
	return
}

// Float returns true if reflect.Value.Float() would succeed.
func Float(rtype r.Type) (ret bool) {
	switch rtype.Kind() {
	case r.Float32, r.Float64:
		ret = true
	}
	return
}

// Number returns true if the passed kind is a float or int
func Number(rtype r.Type) bool {
	return Int(rtype) || Float(rtype)
}

// String kind.
func String(rtype r.Type) bool {
	return rtype.Kind() == r.String
}

// IdentId tests reflect.TypeOf(iden.Id)
func IdentId(rtype r.Type) bool {
	return rtype == identId
}

// Object tests reflect.TypeOf(rt.Object)
func Object(rtype r.Type) bool {
	return rtype == object
}

// BoolEval tests reflect.TypeOf(rt.BoolEval)
func BoolEval(rtype r.Type) bool {
	return rtype.Implements(TypeBoolEval)
}

// NumberEval tests reflect.TypeOf(rt.NumberEval)
func NumberEval(rtype r.Type) bool {
	return rtype.Implements(TypeNumEval)
}

// TextEval tests reflect.TypeOf(rt.TextEval)
func TextEval(rtype r.Type) bool {
	return rtype.Implements(TypeTextEval)
}

// ObjectEval tests reflect.TypeOf(rt.ObjectEval)
func ObjectEval(rtype r.Type) bool {
	return rtype.Implements(TypeObjEval)
}

// NumListEval tests reflect.TypeOf(rt.NumListEval)
func NumListEval(rtype r.Type) bool {
	return rtype.Implements(TypeNumListEval)
}

// TextListEval tests reflect.TypeOf(rt.TextListEval)
func TextListEval(rtype r.Type) bool {
	return rtype.Implements(TypeTextListEval)
}

// ObjListEval tests reflect.TypeOf(rt.ObjListEval)
func ObjListEval(rtype r.Type) bool {
	return rtype.Implements(TypeObjListEval)
}

// switches dont work well with .Interface().(type) when dst is nil.
var object = r.TypeOf((*rt.Object)(nil)).Elem()
var identId = r.TypeOf((*ident.Id)(nil)).Elem()
var TypeBoolEval = r.TypeOf((*rt.BoolEval)(nil)).Elem()
var TypeNumEval = r.TypeOf((*rt.NumberEval)(nil)).Elem()
var TypeTextEval = r.TypeOf((*rt.TextEval)(nil)).Elem()
var TypeObjEval = r.TypeOf((*rt.ObjectEval)(nil)).Elem()
var TypeNumListEval = r.TypeOf((*rt.NumListEval)(nil)).Elem()
var TypeTextListEval = r.TypeOf((*rt.TextListEval)(nil)).Elem()
var TypeObjListEval = r.TypeOf((*rt.ObjListEval)(nil)).Elem()

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
