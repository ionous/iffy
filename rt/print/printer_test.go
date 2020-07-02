package print

import (
	"io"
	"testing"
)

func TestBracket(t *testing.T) {
	var buffer Spanner
	w := Parens(&buffer)
	io.WriteString(w, "hello")
	io.WriteString(w, "you")
	w.Close()
	if str := buffer.String(); str != "( hello you )" {
		t.Fatal("mismatched", str)
	}
}

func TestManualBracket(t *testing.T) {
	var buffer Spanner
	w := &buffer
	io.WriteString(w, "hello")
	io.WriteString(w, "( you )")
	io.WriteString(w, "guys")
	if str := buffer.String(); str != "hello ( you ) guys" {
		t.Fatal("mismatched", str)
	}
}

func TestCapitalize(t *testing.T) {
	var buffer Spanner
	w := Capitalize(&buffer)
	io.WriteString(w, "hello")
	io.WriteString(w, "you")
	if str := buffer.String(); str != "Hello you" {
		t.Fatal("mismatched", str)
	}
}

func TestLowercase(t *testing.T) {
	var buffer Spanner
	w := Lowercase(&buffer)
	io.WriteString(w, "Hello")
	io.WriteString(w, "Hugh")
	if str := buffer.String(); str != "hello hugh" {
		t.Fatal("mismatched", str)
	}
}

func TestTitlecase(t *testing.T) {
	var buffer Spanner
	w := TitleCase(&buffer)
	io.WriteString(w, "hello")
	io.WriteString(w, "you")
	if str := buffer.String(); str != "Hello You" {
		t.Fatal("mismatched", str)
	}
}
