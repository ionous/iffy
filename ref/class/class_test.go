package class

import (
	"github.com/ionous/iffy/ident"
	testify "github.com/stretchr/testify/assert"
	r "reflect"
	"testing"
)

func TestClass(t *testing.T) {
	assert := testify.New(t)

	// BaseClass provides a simple object with every common type.
	type BaseClass struct {
		Name      string `if:"id"`
		MyProp    string
		Overriden string
	}

	// DerivedClass extends BaseClass, and "hides" the property called Overriden.
	type DerivedClass struct {
		BaseClass `if:"parent"`
		Overriden string
	}

	baseClass := r.TypeOf((*BaseClass)(nil)).Elem()
	derivedClass := r.TypeOf((*DerivedClass)(nil)).Elem()

	assert.Equal(ident.IdOf("baseClass"), Id(baseClass))
	assert.Equal("base class", FriendlyName(baseClass))

	assert.Equal(ident.IdOf("$derivedClass"), Id(derivedClass))
	assert.Equal("derived class", FriendlyName(derivedClass))

	if p, ok := Parent(baseClass); assert.False(ok) {
		assert.Nil(p)
	}

	if p, ok := Parent(derivedClass); assert.True(ok) {
		assert.Equal(baseClass, p)
	}

	assert.True(IsCompatible(baseClass, "base class"))
	assert.False(IsCompatible(baseClass, "derived class"))

	assert.True(IsCompatible(derivedClass, "derived class"))
	assert.True(IsCompatible(derivedClass, "base class"))

	assert.Len(PropertyPath(baseClass, "doesnt exist"), 0)
	assert.Len(PropertyPath(baseClass, "my prop"), 1)
	assert.Len(PropertyPath(derivedClass, "my prop"), 2)
	assert.Len(PropertyPath(baseClass, "overriden"), 1)
	assert.Len(PropertyPath(derivedClass, "overriden"), 1)
}
