package core

import (
	"testing"

	"github.com/ionous/iffy/rt/test"
)

func TestMatches(t *testing.T) {
	var run test.PanicRuntime
	// test a valid regexp
	// loop to verify(ish) the cache
	m := &Matches{Text: &Text{"gophergopher"}, Pattern: "(gopher){2}"}
	for i := 0; i < 2; i++ {
		if ok, e := m.GetBool(&run); e != nil {
			t.Fatal(e)
		} else if !ok.Bool() {
			t.Fatal("mismatch")
		}
	}
	// test that a bad expression will fail
	// loop to verify(ish) the cache
	fail := &Matches{Pattern: "("}
	for i := 0; i < 2; i++ {
		if _, e := fail.GetBool(&run); e == nil {
			t.Fatal("expected error")
		} else if _, e := fail.GetBool(&run); e == nil {
			t.Fatal("expected error")
		}
	}
}
