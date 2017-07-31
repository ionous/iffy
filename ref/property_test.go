package ref

import (
	"github.com/ionous/sliceOf"
	testify "github.com/stretchr/testify/assert"
	r "reflect"
	"testing"
)

// TestProperties verifies property generation:
func TestProperties(t *testing.T) {
	type Base struct {
		Index float64
	}
	type Derived struct {
		Base `if:"parent"`
		Text string
	}
	type Embed struct {
		Apples string
		Pears  float64
	}
	type Embedded struct {
		Embed
	}
	type MoreComplex struct {
		Text string
		Embed
		Base  `if:"parent"`
		Float float64
	}
	assert := testify.New(t)
	{
		p, pi, props, e := MakeProperties(r.TypeOf((*Base)(nil)).Elem())
		assert.NoError(e)
		assert.Nil(p)
		assert.Empty(pi)
		assert.Len(props, 1)
	}
	{
		p, pi, props, e := MakeProperties(r.TypeOf((*Derived)(nil)).Elem())
		assert.NoError(e)
		assert.NotNil(p)
		assert.Equal(sliceOf.Int(0), pi)
		assert.Len(props, 1)
	}

	{
		p, pi, props, e := MakeProperties(r.TypeOf((*Embedded)(nil)).Elem())
		assert.NoError(e)
		assert.Nil(p)
		assert.Empty(pi)
		assert.Len(props, 2)
	}
	{
		p, pi, props, e := MakeProperties(r.TypeOf((*MoreComplex)(nil)).Elem())
		assert.NoError(e)
		assert.NotNil(p)
		assert.Equal(sliceOf.Int(2), pi)
		assert.Len(props, 4)
	}
}
