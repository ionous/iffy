package ref

import (
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
		Base
		Text string
	}
	assert := testify.New(t)
	{
		p, pi, props, e := MakeProperties(r.TypeOf((*Base)(nil)).Elem())
		assert.NoError(e)
		assert.Zero(pi)
		assert.Nil(p)
		assert.Len(props, 1)
	}
	{
		p, pi, props, e := MakeProperties(r.TypeOf((*Derived)(nil)).Elem())
		assert.NoError(e)
		assert.Zero(pi)
		assert.NotNil(p)
		assert.Len(props, 1)
	}
}
