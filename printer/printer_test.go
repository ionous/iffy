package printer

import (
	"io"
	"testing"

	testify "github.com/stretchr/testify/assert"
)

func TestBracket(t *testing.T) {
	assert := testify.New(t)
	//
	var buffer Spanner
	w := Bracket(&buffer)
	io.WriteString(w, "hello")
	io.WriteString(w, "you")
	w.Close()
	assert.Equal("( hello you )", buffer.String())
}

func TestManualBracket(t *testing.T) {
	assert := testify.New(t)
	//
	var buffer Spanner
	w := &buffer
	io.WriteString(w, "hello")
	io.WriteString(w, "( you )")
	io.WriteString(w, "guys")
	assert.Equal("hello ( you ) guys", buffer.String())
}

func TestCapitalize(t *testing.T) {
	assert := testify.New(t)
	//
	var buffer Spanner
	w := Capitalize(&buffer)
	io.WriteString(w, "hello")
	io.WriteString(w, "you")
	assert.Equal("Hello you", buffer.String())
}

func TestLowercase(t *testing.T) {
	assert := testify.New(t)
	//
	var buffer Spanner
	w := Lowercase(&buffer)
	io.WriteString(w, "Hello")
	io.WriteString(w, "Hugh")
	assert.Equal("hello hugh", buffer.String())
}

func TestTitlecase(t *testing.T) {
	assert := testify.New(t)
	//
	var buffer Spanner
	w := TitleCase(&buffer)
	io.WriteString(w, "hello")
	io.WriteString(w, "you")
	assert.Equal("Hello You", buffer.String())
}
