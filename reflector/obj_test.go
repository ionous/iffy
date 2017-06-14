package reflector

import (
	"github.com/ionous/iffy/ref"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestObjects(t *testing.T) {
	assert := assert.New(t)
	test := func(n ref.Object) {
		var name string
		assert.NoError(n.GetValue("name", &name))

		// test text
		{
			var v string
			assert.NoError(n.GetValue("text", &v))
			assert.Empty(v)
			assert.NoError(n.SetValue("text", "something"))
			assert.NoError(n.GetValue("text", &v))
			assert.Equal("something", v)
		}

		// text number: float and int

		// test state: bool

		// pointers.... and relations?
	}
	b := BaseClass{Name: "first-instance"}
	d := DerivedClass{BaseClass{Name: "second-instance"}}
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
