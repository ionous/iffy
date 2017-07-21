package printer

import (
	testify "github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestPrintSep(t *testing.T) {
	assert := testify.New(t)
	//
	if s, e := write(AndSeparator, "pizza"); assert.NoError(e) {
		assert.Equal("pizza", s)
	}
	if s, e := write(AndSeparator, "apple", "hedgehog", "washington", "mushroom"); assert.NoError(e) {
		assert.Equal("apple, hedgehog, washington, and mushroom", s)
	}
	//
	if s, e := write(OrSeparator, "pistachio"); assert.NoError(e) {
		assert.Equal("pistachio", s)
	}
	if s, e := write(OrSeparator, "apple", "hedgehog", "washington", "mushroom"); assert.NoError(e) {
		assert.Equal("apple, hedgehog, washington, or mushroom", s)
	}
}

func write(sep func(w io.Writer) *Sep, names ...string) (ret string, err error) {
	var buffer Span
	w := sep(&buffer)
	for _, n := range names {
		if _, e := io.WriteString(w, n); e != nil {
			err = e
			break
		}
	}
	if err == nil {
		// normally PopWriter would call close, but we arent using the runtime here.
		if e := w.Close(); e != nil {
			err = e
		} else {
			ret = buffer.String()
		}
	}
	return
}
