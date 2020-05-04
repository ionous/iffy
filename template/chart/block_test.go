package chart

import (
	"testing"

	"github.com/ionous/errutil"
)

func TestBlocks(t *testing.T) {
	if e := testBlock(t, "", ""); e != nil {
		t.Fatal(e)
	}
	if e := testBlock(t, "abc", "abc"); e != nil {
		t.Fatal(e)
	}
	if e := testBlock(t, "{}", "{}"); e != nil {
		t.Fatal(e)
	}
	// mixed: front, end
	if e := testBlock(t, "abc{}", "abc{}"); e != nil {
		t.Fatal(e)
	}
	if e := testBlock(t, "{}abc", "{}abc"); e != nil {
		t.Fatal(e)
	}
	if e := testBlock(t, "{}{}", "{}{}"); e != nil {
		t.Fatal(e)
	}
	// long
	if e := testBlock(t, "abc{}d{}efg{}z", "abc{}d{}efg{}z"); e != nil {
		t.Fatal(e)
	}
}

func TestKeys(t *testing.T) {
	if e := testBlock(t, "{key}", "{key:}"); e != nil {
		t.Fatal(e)
	}
	if e := testBlock(t, "{ key }", "{key:}"); e != nil {
		t.Fatal(e)
	}
	if e := testBlock(t, "{1}", ignoreResult); e == nil {
		t.Fatal(e)
	}
	if e := testBlock(t, "{key1}", ignoreResult); e == nil {
		t.Fatal(e)
	}
}

func TestTrim(t *testing.T) {
	if e := testBlock(t, "{~~}", "{}"); e != nil {
		t.Fatal(e)
	}
	if e := testBlock(t, "    {~~}    ", "{}"); e != nil {
		t.Fatal(e)
	}

	if e := testBlock(t, "abc{~ }", "abc{}"); e != nil {
		t.Fatal(e)
	}
	if e := testBlock(t, "abc   {~ }", "abc{}"); e != nil {
		t.Fatal(e)
	}

	if e := testBlock(t, "{ ~}abc", "{}abc"); e != nil {
		t.Fatal(e)
	}
	if e := testBlock(t, "{ ~}    abc", "{}abc"); e != nil {
		t.Fatal(e)
	}

	if e := testBlock(t, "{ ~}{~ }", "{}{}"); e != nil {
		t.Fatal(e)
	}
	if e := testBlock(t, "{ ~}  {~ }", "{}{}"); e != nil {
		t.Fatal(e)
	}
	if e := testBlock(t, "abc {  }  d {   } efg  {  }z", "abc {}  d {} efg  {}z"); e != nil {
		t.Fatal(e)
	}
	if e := testBlock(t, "abc {~ }  d {~ ~} efg  {~ }z", "abc{}  d{}efg{}z"); e != nil {
		t.Fatal(e)
	}
}

func TestKeyTrim(t *testing.T) {
	if e := testBlock(t, "  {~key}", "{key:}"); e != nil {
		t.Fatal(e)
	}
	if e := testBlock(t, "  { key}", "  {key:}"); e != nil {
		t.Fatal(e)
	}
	if e := testBlock(t, "  {~key~}  ", "{key:}"); e != nil {
		t.Fatal(e)
	}
	if e := testBlock(t, "  {key~}  ", "  {key:}"); e != nil {
		t.Fatal(e)
	}
	if e := testBlock(t, "  {key}  ", "  {key:}  "); e != nil {
		t.Fatal(e)
	}
}

func testBlock(t *testing.T, str string, want string) (err error) {
	t.Log("test:", str)
	p := BlockParser{factory: EmptyFactory{}}
	if e := Parse(&p, str); e != nil {
		err = e
	} else if res, e := p.GetDirectives(); e != nil {
		err = e
	} else if want != ignoreResult {
		got := Format(res)
		if got == want {
			t.Log("ok:", got)
		} else {
			err = mismatched(want, got)
		}
	}
	return
}

func mismatched(want, got string) error {
	return errutil.Fmt("want(%d): %s; != got(%d): %s.", len(want), want, len(got), got)
}
