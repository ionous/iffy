package ref

import (
	"github.com/ionous/iffy/ref/unique"
	. "github.com/ionous/iffy/tests"
	"github.com/stretchr/testify/suite"
	r "reflect"
	"testing"
)

//
func TestClassSuite(t *testing.T) {
	suite.Run(t, new(ClassSuite))
}

type ClassSuite struct {
	suite.Suite
}

// type Expected struct {
// 	name string
// 	kind rt.PropertyType
// }

// func expected() []Expected {
// 	return []Expected{
// 		{"Name", rt.Text},
// 		{"Num", rt.Number},
// 		{"Text", rt.Text},
// 		{"Object", rt.Pointer},
// 		{"Nums", rt.Number | rt.Array},
// 		{"Texts", rt.Text | rt.Array},
// 		{"Objects", rt.Pointer | rt.Array},
// 		{"State", rt.State},
// 		{"Labeled", rt.State},
// 	}
// }

func (assert *ClassSuite) TestClass() {
	//
	cs := NewClasses()
	baseType := r.TypeOf((*BaseClass)(nil)).Elem()
	derivedType := r.TypeOf((*DerivedClass)(nil)).Elem()

	// id field tests
	if path, ok := unique.PathOf(baseType, "id"); assert.True(ok) {
		field := baseType.FieldByIndex(path)
		assert.Equal(field.Name, "Name")
	}
	if path, ok := unique.PathOf(derivedType, "id"); assert.True(ok) {
		field := derivedType.FieldByIndex(path)
		assert.Equal(field.Name, "Name")
	}

	// add and retrieve base class:
	var baseClass RefClass
	if ref, e := cs.RegisterClass(baseType); assert.NoError(e) {
		baseClass = ref
		// base class tests:
		assert.Equal("$baseClass", baseClass.GetId())
		_, parentExists := baseClass.GetParent()
		assert.False(parentExists)

		// // test the property interfaces:
		// assert.Equal(len(expected()), baseClass.NumProperty())
		// first := baseClass.PropertyNum(0)
		// if name, ok := baseClass.GetProperty("name"); assert.True(ok) {
		// 	assert.Equal(first, name)
		// }
		// // test find by choice:
		// if state, ok := baseClass.GetProperty("state"); assert.True(ok) {
		// 	if p, ok := baseClass.GetPropertyByChoice("yes"); assert.True(ok) {
		// 		assert.Equal(state, p)
		// 	}
		// }
		// derived class tests:
		if derived, e := cs.RegisterClass(derivedType); assert.NoError(e) {
			assert.Equal("$derivedClass", derived.GetId())
			if p, ok := derived.GetParent(); assert.True(ok) {
				assert.Equal(baseClass, p)
				assert.True(p.IsCompatible(p.GetId()))
			}
		}
		// class set verification:
		assert.Contains(cs.ClassMap, "$baseClass")
		assert.Contains(cs.ClassMap, "$derivedClass")
	}
}

// func (assert *ClassSuite) TestClassProperties() {
// 	baseType := r.TypeOf((*BaseClass)(nil)).Elem()
// 	if parent, _, props, e := MakeProperties(baseType); assert.NoError(e) {
// 		assert.Nil(parent)
// 		for i, v := range expected() {
// 			p := props[i]
// 			assert.Equal(id.MakeId(v.name), p.GetId(), v.name)
// 			assert.Equal(v.kind, p.GetType(), v.name)
// 		}
// 		// var md unique.Metadata
// 		// if assert.Len(md, 2) {
// 		// 	assert.Equal("Name", md["id"])
// 		// 	assert.Equal("baseType classes", md["plural"])
// 		// }
// 	}
// 	derived := r.TypeOf((*DerivedClass)(nil)).Elem()
// 	if parent, _, props, e := MakeProperties(derived); assert.NoError(e) {
// 		assert.Equal(baseType, parent)
// 		assert.Len(props, 0)
// 	}
// 	// var dd unique.Metadata
// 	// if assert.Len(dd, 1) {
// 	// 	assert.Equal("derives", dd["plural"])
// 	// }
// }
