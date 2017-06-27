package reflector

import (
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/reflector/unique"
	"github.com/ionous/iffy/rt"
	. "github.com/ionous/iffy/tests"
	"github.com/stretchr/testify/assert"
	r "reflect"
	"testing"
)

type BaseClass struct {
	Name    string `if:"id,plural:base classes"`
	Num     float64
	Text    string
	Object  *BaseClass
	Nums    []float64
	Texts   []string
	Objects []*BaseClass
	State   TriState
	Labeled bool
}

type DerivedClass struct {
	BaseClass `if:"plural:derives"`
}

type Expected struct {
	name string
	kind rt.PropertyType
}

func expected() []Expected {
	return []Expected{
		{"Name", rt.Text},
		{"Num", rt.Number},
		{"Text", rt.Text},
		{"Object", rt.Pointer},
		{"Nums", rt.Number | rt.Array},
		{"Texts", rt.Text | rt.Array},
		{"Objects", rt.Pointer | rt.Array},
		{"State", rt.State},
		{"Labeled", rt.State},
	}
}

func TestClass(t *testing.T) {
	assert := assert.New(t)
	//
	cs := make(Classes)
	base := r.TypeOf((*BaseClass)(nil)).Elem()
	// add and retrieve base class:
	var baseClass *RefClass
	if ref, e := cs.addClass(base); assert.NoError(e) {
		baseClass = ref
		// base class tests:
		assert.Equal("$baseClass", baseClass.GetId())
		_, parentExists := baseClass.GetParent()
		assert.False(parentExists)
		// id field
		assert.Equal("Name", baseClass.findId())
		// test the property interfaces:
		assert.Equal(len(expected()), baseClass.NumProperty())
		first := baseClass.PropertyNum(0)
		if name, ok := baseClass.GetProperty("name"); assert.True(ok) {
			assert.Equal(first, name)
		}
		// test find by choice:
		if state, ok := baseClass.GetProperty("state"); assert.True(ok) {
			if p, ok := baseClass.GetPropertyByChoice("yes"); assert.True(ok) {
				assert.Equal(state, p)
			}
		}
		// derived class tests:
		derivedType := r.TypeOf((*DerivedClass)(nil)).Elem()
		if derived, e := cs.addClass(derivedType); assert.NoError(e) {
			assert.Equal("$derivedClass", derived.GetId())
			if p, ok := derived.GetParent(); assert.True(ok) {
				assert.Equal(baseClass, p)
				assert.True(p.IsCompatible(p.GetId()))
			}
			// id field
			assert.Equal("Name", derived.findId())
		}
		// class set verification:
		assert.Contains(cs, "$baseClass")
		assert.Contains(cs, "$derivedClass")
	}
}

func TestClassProperties(t *testing.T) {
	assert := assert.New(t)
	base := r.TypeOf((*BaseClass)(nil)).Elem()
	var md unique.Metadata
	if parent, _, props, e := MakeProperties(base, &md); assert.NoError(e) {
		assert.Nil(parent)
		for i, v := range expected() {
			p := props[i]
			assert.Equal(id.MakeId(v.name), p.GetId(), v.name)
			assert.Equal(v.kind, p.GetType(), v.name)
		}
		if assert.Len(md, 2) {
			assert.Equal("Name", md["id"])
			assert.Equal("base classes", md["plural"])
		}
	}
	var dd unique.Metadata
	derived := r.TypeOf((*DerivedClass)(nil)).Elem()
	if parent, _, props, e := MakeProperties(derived, &dd); assert.NoError(e) {
		assert.Equal(base, parent)
		assert.Len(props, 0)
	}
	if assert.Len(dd, 1) {
		assert.Equal("derives", dd["plural"])
	}
}
