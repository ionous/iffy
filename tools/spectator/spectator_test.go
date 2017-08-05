package main

// import (
// 	"github.com/ionous/errutil"
// 	"github.com/ionous/sliceOf"
// 	testify "github.com/stretchr/testify/assert"
// 	"strings"
// 	"testing"
// )

// func TestInput(t *testing.T) {
// 	assert := testify.New(t)
// 	{
// 		x, e := split(`NumWorder "NumWord A B Last"`)
// 		assert.NoError(e)
// 		assert.Equal(sliceOf.String("NumWorder", "NumWord", "A", "B", "Last"), x)
// 	}
// 	{
// 		x, e := split(`NumWorder "NumWord"`)
// 		assert.NoError(e)
// 		assert.Equal(sliceOf.String("NumWorder", "NumWord"), x)
// 	}
// 	{
// 		_, e := split(`NumWorder`)
// 		assert.Error(e)
// 	}
// }
