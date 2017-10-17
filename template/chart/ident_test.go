package chart

import (
	"testing"
)

func TestIdent(t *testing.T) {
	// returns point of failure
	test := func(str string) (ret string, err error) {
		var p identParser
		if end := parse(&p, str); end > 0 {
			err = endpointError(end)
		} else if v, e := p.GetName(); e != nil {
			err = e
		} else {
			ret = v
		}
		return
	}
	fails := func(str string) {
		if v, e := test(str); e != nil {
			t.Log(str, "ok", e)
		} else {
			t.Fatal(str, "expected error", v)
		}
	}
	succeeds := func(str string) {
		if res, e := test(str); e != nil {
			t.Fatal(e, "for:", str)
		} else if str != res {
			t.Fatal(str, "unexpected value", res)

		}
	}
	fails("")
	fails(" a")
	fails("b ")
	fails("1a")
	succeeds("a1")
	succeeds("abc")
}
