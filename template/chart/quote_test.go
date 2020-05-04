package chart

import (
	"testing"
)

func TestQuotes(t *testing.T) {
	if e, x := testQ(t, "'singles'", "singles"); e != nil {
		t.Fatal(e, x)
	}
	if e, x := testQ(t, `"doubles"`, "doubles"); e != nil {
		t.Fatal(e, x)
	}
	if e, x := testQ(t, "'escape\"'", "escape\""); e != nil {
		t.Fatal(e, x)
	}
	if e, x := testQ(t, `"\\"`, "\\"); e != nil {
		t.Fatal(e, x)
	}
	if e, x := testQ(t, string([]rune{'"', '\\', 'a', '"'}), "\a"); e != nil {
		t.Fatal(e, x)
	}
	if e, _ := testQ(t, string([]rune{'"', '\\', 'g', '"'}),
		ignoreResult); e == nil {
		t.Fatal(e)
	}
}

func testQ(t *testing.T, str, want string) (err error, ret interface{}) {
	t.Log("test:", str)
	var p QuoteParser
	if e := Parse(&p, str); e != nil {
		err = e
	} else if got, e := p.GetString(); e != nil {
		err = e
	} else if want != ignoreResult {
		if got != want {
			err = mismatched(want, got)
		} else {
			t.Log("ok", got)
		}
	}
	return err, str
}
