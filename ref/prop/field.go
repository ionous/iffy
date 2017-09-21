package prop

import (
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/coerce"
	r "reflect"
)

// Field of an object containing a simple value, or an array of simple values.
type Field struct {
	field r.StructField
	value r.Value // container of Value property
}

// MakeField constructor. panics if not value.CanSet().
func MakeField(field r.StructField, value r.Value) Field {
	return Field{field, value}
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

// Value returns the current value of the property.
func (p Field) Value() r.Value {
	return p.value
}

// SetValue to change the current value of the property.
// Returns an error if the passed value is not compatible with Type().
func (p Field) SetValue(src r.Value) (err error) {
	return coerce.Value(p.value, src)
}
