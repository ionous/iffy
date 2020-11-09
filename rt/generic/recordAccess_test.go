package generic

import (
	"testing"

	"github.com/ionous/iffy/object"
)

//
func TestRecordAccess(t *testing.T) {
	q := newRecordAccessTest()
	//
	t.Run("numbers", func(t *testing.T) {
		q := newRecordAccessTest()
		// beep, number
		if beep, e := q.GetField(object.Variables, "beep"); e != nil {
			t.Fatal(e)
		} else if el, e := beep.GetNamedField("d"); e != nil {
			t.Fatal(e)
		} else if v, e := el.GetNumber(); e != nil {
			t.Fatal(e)
		} else if v != 0 {
			t.Fatal("not default", v)
		} else if e := beep.SetNamedField("d", FloatOf(5)); e != nil {
			t.Fatal(e)
		} else if el, e := beep.GetNamedField("d"); e != nil {
			t.Fatal(e)
		} else if v, e := el.GetNumber(); e != nil {
			t.Fatal(e)
		} else if v != 5 {
			t.Fatal("not changed", v)
		}
	})
	t.Run("text", func(t *testing.T) {
		if boop, e := q.GetField(object.Variables, "boop"); e != nil {
			t.Fatal(e)
		} else if el, e := boop.GetNamedField("t"); e != nil {
			t.Fatal(e)
		} else if v, e := el.GetText(); e != nil {
			t.Fatal(e)
		} else if len(v) > 0 {
			t.Fatal("not default", v)
		} else if e := boop.SetNamedField("t", StringOf("xyzzy")); e != nil {
			t.Fatal(e)
		} else if el, e := boop.GetNamedField("t"); e != nil {
			t.Fatal(e)
		} else if v, e := el.GetText(); e != nil {
			t.Fatal(e)
		} else if v != "xyzzy" {
			t.Fatal("not changed", v)
		}
	})
	t.Run("aspects", func(t *testing.T) {
		if beep, e := q.GetField(object.Variables, "beep"); e != nil {
			t.Fatal(e)
		} else if el, e := beep.GetNamedField("x"); e != nil {
			t.Fatal(e)
		} else if v, e := el.GetBool(); e != nil {
			t.Fatal(e)
		} else if v != false {
			t.Fatal("not default", v)
		} else if e := beep.SetNamedField("x", BoolOf(true)); e != nil {
			t.Fatal(e)
		} else if el, e := beep.GetNamedField("x"); e != nil {
			t.Fatal(e)
		} else if v, e := el.GetBool(); e != nil {
			t.Fatal(e)
		} else if v != true {
			t.Fatal("not changed", v)
		} else if el, e := beep.GetNamedField("a"); e != nil {
			t.Fatal(e)
		} else if v, e := el.GetText(); e != nil {
			t.Fatal(e)
		} else if v != "x" {
			t.Fatal(e)
		} else if e := beep.SetNamedField("a", StringOf("w")); e != nil {
			t.Fatal(e)
		} else if el, e := beep.GetNamedField("w"); e != nil {
			t.Fatal(e)
		} else if v, e := el.GetBool(); e != nil {
			t.Fatal(e)
		} else if v != true {
			t.Fatal("aspect not changed")
		}
	})
	t.Run("failures", func(t *testing.T) {
		if beep, e := q.GetField(object.Variables, "beep"); e != nil {
			t.Fatal(e)
		} else if _, e := beep.GetNamedField("nope"); e == nil {
			t.Fatal("expected no such field")
		} else if e := beep.SetNamedField("a", True); e == nil {
			t.Fatal("aspects should be set with strings")
		} else if e := beep.SetNamedField("x", Empty); e == nil {
			t.Fatal("traits should be set with bools")
		} else if e := beep.SetNamedField("x", False); e == nil {
			// we dont have support for opposite values right now.
			t.Fatal("traits should be set with true values only")
		} else if _, e := q.GetField(object.Variables, "blip"); e == nil {
			// really this tests the testing mock not the record code, but whatev.
			t.Fatal("expected no such variable")
		}
	})

}
