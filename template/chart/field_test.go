package chart

import (
	"strings"
	"testing"

	"github.com/kr/pretty"
)

func TestFields(t *testing.T) {
	// fieldFails(t, "a.")
	// fieldFails(t, "b")
	// fieldFails(t, "")
	fieldSucceeds(t, ".a")
	fieldSucceeds(t, ".a.b")
	fieldSucceeds(t, ".a.b1.c2")
}

func testField(t *testing.T, str string) (ret []string, err error) {
	t.Logf("parsing %q", str)
	var p FieldParser
	if e := Parse(&p, str); e != nil {
		err = e
	} else if n, e := p.GetFields(); e != nil {
		err = e
	} else {
		ret = n
	}
	return
}

func fieldFails(t *testing.T, str string) {
	if v, e := testField(t, str); e != nil {
		t.Log("ok, error:", e)
	} else {
		t.Fatal("expected error", v)
	}
}

func fieldSucceeds(t *testing.T, str string) {
	if res, e := testField(t, str); e != nil {
		t.Fatal(e, "for:", str)
	} else {
		var match []string
		if len(str) > 0 {
			match = strings.Split(str[1:], ".")
		}
		diff := pretty.Diff(res, match)
		if len(diff) > 0 {
			t.Fatal("unexpected value", res)
		}
	}
}
