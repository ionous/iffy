package sds

import "testing"

func TestNote(t *testing.T) {
	{
		var a Note = "chicken::egg@id"
		if f, n, d := a.Extract(); f != "chicken" || n != "egg" || d != "id" {
			t.Fatal("mismatch", f, n, d)
		}
	}
	{
		var a Note = "::egg@id"
		if f, n, d := a.Extract(); f != "egg" || n != "egg" || d != "id" {
			t.Fatal("mismatch", f, n, d)
		}
	}
	{
		var a Note = "::egg"
		if f, n, d := a.Extract(); f != "egg" || n != "egg" || d != "" {
			t.Fatal("mismatch", f, n, d)
		}
	}
	{
		var a Note = "::@id"
		if f, n, d := a.Extract(); f != "" || n != "" || d != "id" {
			t.Fatal("mismatch", f, n, d)
		}
	}
	{
		var a Note = "chicken::@id"
		if f, n, d := a.Extract(); f != "chicken" || n != "" || d != "id" {
			t.Fatal("mismatch", f, n, d)
		}
	}
	{
		var a Note = "chicken::egg"
		if f, n, d := a.Extract(); f != "chicken" || n != "egg" || d != "" {
			t.Fatal("mismatch", f, n, d)
		}
	}
	{
		var a Note = "chicken"
		if f, n, d := a.Extract(); f != "chicken" || n != "" || d != "" {
			t.Fatal("mismatch", f, n, d)
		}
	}
}
