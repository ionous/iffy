package rt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPropertyType(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(NumberArray, Number|Array)
	assert.Equal(TextArray, Text|Array)
	assert.Equal(PointerArray, Pointer|Array)
}
