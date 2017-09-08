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

	str := Field{obj, obj.Type().Field(0)}
	num := Field{obj, obj.Type().Field(1)}
	tri := Field{obj, obj.Type().Field(2)}
	yes := State{tri, 1}
	no := State{tri, 0}
	//
	// var _ rt.Property = str
	// var _ rt.Property = yes
	//
	assert := assert.New(t)
	// test that we can see expected values via the interface
	assert.Equal("foo", str.Value())
	assert.Equal(5, num.Value())
	// test that we can coerce via the value interface
	// -- TestCoerce runs the full gamut of coercion tests.
	num.SetValue(3.2)
	assert.Equal(3, num.Value())
	// test that changing one state changes them all
	assert.Equal(true, yes.Value())
	assert.Equal(false, no.Value())
	no.SetValue(true)
	assert.Equal(false, yes.Value())
	assert.Equal(true, no.Value())
}
