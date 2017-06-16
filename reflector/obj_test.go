package reflector

import (
	"github.com/ionous/iffy/ref"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestObjectGetSet(t *testing.T) {
	assert := assert.New(t)
	test := func(n ref.Object) {
		state := false
		if e := n.GetValue("yes", &state); assert.NoError(e) {
			assert.True(state)
		}
		if e := n.GetValue("maybe", &state); assert.NoError(e) {
			assert.False(state)
		}

		// var name string
		// assert.NoError(n.GetValue("name", &name))

		// test text
		// {
		// 	var v string
		// 	assert.NoError(n.GetValue("text", &v))
		// 	assert.Empty(v)
		// 	assert.NoError(n.SetValue("text", "something"))
		// 	assert.NoError(n.GetValue("text", &v))
		// 	assert.Equal("something", v)
		// }

		// text number: float and int

		// test state: bool

		// pointers.... and relations?
	}
	b := BaseClass{Name: "first-instance", State: Yes}
	d := DerivedClass{BaseClass{Name: "second-instance", State: Yes}}
	if m, e := MakeModel(b, d); assert.NoError(e) {
		if n, ok := m.GetObject("first-instance"); assert.True(ok) {
			assert.Equal("$firstInstance", n.GetId())
			cls := n.GetClass()
			assert.NotNil(cls)
			assert.Equal("$baseClass", cls.GetId())
			//
			test(n)
		}
		if d, ok := m.GetObject("second-instance"); assert.True(ok) {
			assert.Equal("$secondInstance", d.GetId())
			cls := d.GetClass()
			assert.NotNil(cls)
			assert.Equal("$derivedClass", cls.GetId())
			//
			test(d)
		}
	}
}
