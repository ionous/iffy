package reflector

import (
	"github.com/stretchr/testify/assert"
	r "reflect"
	"testing"
)

func TestMetadata(t *testing.T) {
	assert := assert.New(t)
	//`if:"id,plural:base classes"`
	MakeMetadata := func(tag string) (out Metadata) {
		s := r.StructTag(tag).Get("if")
		if len(s) > 0 {
			out = make(Metadata)
			out.AddString(s, "id")
		}
		return
	}
	{
		m := MakeMetadata(`if:"id"`)
		assert.Len(m, 1)
		t.Log("WWEWLK", m)
		assert.Equal("id", m["id"])
	}
	{
		m := MakeMetadata(`if:"plural:tests"`)
		assert.Len(m, 1)
		assert.Equal("tests", m["plural"])
	}
	{
		m := MakeMetadata(`if:"id,plural:tests"`)
		assert.Len(m, 2)
		assert.Equal("id", m["id"])
		assert.Equal("tests", m["plural"])
	}
	{
		m := MakeMetadata(`if:"edge:test:tests"`)
		assert.Len(m, 1)
		assert.Equal("test:tests", m["edge"])
	}
}
