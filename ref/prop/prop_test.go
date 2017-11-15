package prop

import (
	// "github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/tests"
	"github.com/stretchr/testify/assert"
	r "reflect"
	"testing"
)

func TestProp(t *testing.T) {
	type PropClass struct {
		Str string
		Num int
		Tri tests.TriState
	}
	parent := &PropClass{"foo", 5, tests.Yes}
	obj := r.ValueOf(parent).Elem()

	str := Field{obj.Type().Field(0), obj.Field(0)}
	num := Field{obj.Type().Field(1), obj.Field(1)}
	tri := Field{obj.Type().Field(2), obj.Field(2)}
	yes := State{tri, 1}
	no := State{tri, 0}
	//
	// var _ rt.Property = str
	// var _ rt.Property = yes
	//
	assert := assert.New(t)
	// test that we can see expected values via the interface
	assert.Equal("foo", str.Get())
	assert.Equal(5, num.Get())
	// test that we can coerce via the value interface
	// -- TestCoerce runs the full gamut of coercion tests.
	num.Set(3.2)
	assert.Equal(3, num.Get())
	// test that changing one state changes them all
	assert.Equal(true, yes.Get())
	assert.Equal(false, no.Get())
	no.Set(true)
	assert.Equal(false, yes.Get())
	assert.Equal(true, no.Get())
}

func (f Field) Set(i interface{}) {
	f.SetValue(r.ValueOf(i))
}
func (f Field) Get() interface{} {
	return f.Value().Interface()
}
func (s State) Set(i interface{}) {
	s.SetValue(r.ValueOf(i))
}
func (s State) Get() interface{} {
	return s.Value().Interface()
}
