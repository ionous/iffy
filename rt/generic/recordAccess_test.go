package generic

import (
	"testing"

	"github.com/ionous/iffy/object"
)

//
func TestRecordAccess(t *testing.T) {
	q := newRecordAccessTest()
	// beep, number
	if beep, e := q.GetField(object.Variables, "beep"); e != nil {
		t.Fatal(e)
	} else if el, e := beep.GetField("d"); e != nil {
		t.Fatal(e)
	} else if v, e := el.GetNumber(); e != nil {
		t.Fatal(e)
	} else if v != 0 {
		t.Fatal("not default", v)
	} else if e := beep.SetField("d", &Int{Value: 5}); e != nil {
		t.Fatal(e)
	} else if el, e := beep.GetField("d"); e != nil {
		t.Fatal(e)
	} else if v, e := el.GetNumber(); e != nil {
		t.Fatal(e)
	} else if v != 5 {
		t.Fatal("not changed", v)
	}
	// boop, text
	if boop, e := q.GetField(object.Variables, "boop"); e != nil {
		t.Fatal(e)
	} else if el, e := boop.GetField("t"); e != nil {
		t.Fatal(e)
	} else if v, e := el.GetText(); e != nil {
		t.Fatal(e)
	} else if len(v) > 0 {
		t.Fatal("not default", v)
	} else if e := boop.SetField("t", &String{Value: "xyzzy"}); e != nil {
		t.Fatal(e)
	} else if el, e := boop.GetField("t"); e != nil {
		t.Fatal(e)
	} else if v, e := el.GetText(); e != nil {
		t.Fatal(e)
	} else if v != "xyzzy" {
		t.Fatal("not changed", v)
	}
	// beep, trait and aspect
	if beep, e := q.GetField(object.Variables, "beep"); e != nil {
		t.Fatal(e)
	} else if el, e := beep.GetField("x"); e != nil {
		t.Fatal(e)
	} else if v, e := el.GetBool(); e != nil {
		t.Fatal(e)
	} else if v != false {
		t.Fatal("not default", v)
	} else if e := beep.SetField("x", &Bool{Value: true}); e != nil {
		t.Fatal(e)
	} else if el, e := beep.GetField("x"); e != nil {
		t.Fatal(e)
	} else if v, e := el.GetBool(); e != nil {
		t.Fatal(e)
	} else if v != true {
		t.Fatal("not changed", v)
	} else if el, e := beep.GetField("a"); e != nil {
		t.Fatal(e)
	} else if v, e := el.GetText(); e != nil {
		t.Fatal(e)
	} else if v != "x" {
		t.Fatal(e)
	} else if e := beep.SetField("a", &String{Value: "w"}); e != nil {
		t.Fatal(e)
	} else if el, e := beep.GetField("w"); e != nil {
		t.Fatal(e)
	} else if v, e := el.GetBool(); e != nil {
		t.Fatal(e)
	} else if v != true {
		t.Fatal("aspect not changed")
	}
}
