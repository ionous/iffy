package chart

import (
	"github.com/kr/pretty"
	"strings"
	"testing"
)

func TestFields(t *testing.T) {
	test := func(str string) (ret []string, err error) {
		var p fieldParser
		if end := parse(&p, str); end > 0 {
			err = endpointError(end)
		} else if n, e := p.GetFields(); e != nil {
			err = e
		} else {
			ret = n
		}
		return
	}
	fails := func(str string) {
		if v, e := test(str); e != nil {
			t.Log("ok", e)
		} else {
			t.Fatal("expected error", v)
		}
	}
	succeeds := func(str string) {
		if fields, e := test(str); e != nil {
			t.Fatal(e, "for:", str)
		} else {
			diff := pretty.Diff(fields, strings.Split(str, "."))
			if len(diff) > 0 {
				t.Fatal("unexpected value", diff)
			}
		}
	}
	fails("")
	fails("a.")
	fails(".b")
	succeeds("a")
	succeeds("a.b")
	succeeds("a.b1.c2")
}
