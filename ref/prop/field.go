package prop

import (
	"github.com/ionous/iffy/ident"
	r "reflect"
)

// Field of an object containing a simple value, or an array of simple values.
type Field struct {
	parent r.Value // container of fieldValue property
	field  r.StructField
}

func (p Field) fieldValue() r.Value {
	return p.parent.FieldByIndex(p.field.Index)
}

// Id semi-unique identifier of the property slot in its parent object.
// ( In a some given class hierarchy, properties from a parent class may get hidden by a child class. )
func (p Field) Id() ident.Id {
	return ident.IdOf(p.field.Name)
}

// Type of the property slot ( not the type of the value held by the property. )
func (p Field) Type() r.Type {
	return p.field.Type
}

// String name of the property slot.
func (p Field) String() string {
	return p.field.Name
}

// Value returns a snapshot of the current value of the property.
func (p Field) Value() interface{} {
	v := p.fieldValue()
	return v.Interface()
}

// SetValue to change the current value of the property.
// Returns an error if the passed value is not compatible with Type().
func (p Field) SetValue(v interface{}) error {
	dst, src := p.fieldValue(), r.ValueOf(v)
	return CoerceValue(dst, src)
}
