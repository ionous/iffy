package safe_test

import (
	"testing"

	"github.com/ionous/iffy/rt/safe"
	"github.com/ionous/iffy/rt/test"
)

func TestSafety(t *testing.T) {
	var run test.PanicRuntime
	switch e := safe.RunAll(&run, nil); e.(type) {
	case nil:
		t.Log("okay nothing run")
	default:
		t.Fatal(e)
	}
	switch e := safe.Run(&run, nil); e.(type) {
	case safe.MissingEval:
		t.Log("okay", e)
	default:
		t.Fatal(e)
	}
	switch _, e := safe.GetBool(&run, nil); e.(type) {
	case safe.MissingEval:
		t.Log("okay", e)
	default:
		t.Fatal(e)
	}
	switch _, e := safe.GetNumber(&run, nil); e.(type) {
	case safe.MissingEval:
		t.Log("okay", e)
	default:
		t.Fatal(e)
	}
	switch _, e := safe.GetText(&run, nil); e.(type) {
	case safe.MissingEval:
		t.Log("okay", e)
	default:
		t.Fatal(e)
	}
}
