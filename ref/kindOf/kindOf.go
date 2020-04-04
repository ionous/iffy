package kindOf

import (
	r "reflect"

	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/rt"
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

// Interface returns true if the passed type is an interface
func Interface(rtype r.Type) bool {
	return rtype.Kind() == r.Interface
}

// Float returns true if reflect.Value.Float() would succeed.
func Float(rtype r.Type) (ret bool) {
	switch rtype.Kind() {
	case r.Float32, r.Float64:
		ret = true
	}
	return
}

// Number returns true if the passed type is a float or int
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

// NumListEval tests reflect.TypeOf(rt.NumListEval)
func NumListEval(rtype r.Type) bool {
	return rtype.Implements(TypeNumListEval)
}

// TextListEval tests reflect.TypeOf(rt.TextListEval)
func TextListEval(rtype r.Type) bool {
	return rtype.Implements(TypeTextListEval)
}

func Execute(rtype r.Type) bool {
	return rtype.Implements(TypeExecute)
}

// switches dont work well with .Interface().(type) when dst is nil.
var identId = r.TypeOf((*ident.Id)(nil)).Elem()

var TypeString = r.TypeOf("")
var TypeExecute = r.TypeOf((*rt.Execute)(nil)).Elem()
var TypeBoolEval = r.TypeOf((*rt.BoolEval)(nil)).Elem()
var TypeNumEval = r.TypeOf((*rt.NumberEval)(nil)).Elem()
var TypeTextEval = r.TypeOf((*rt.TextEval)(nil)).Elem()
var TypeNumListEval = r.TypeOf((*rt.NumListEval)(nil)).Elem()
var TypeTextListEval = r.TypeOf((*rt.TextListEval)(nil)).Elem()

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
