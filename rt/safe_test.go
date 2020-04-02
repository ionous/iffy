package rt

import "testing"

func TestSafety(t *testing.T) {
	var run Panic
	switch e := Run(&run, nil); e.(type) {
	case nil:
		t.Log("okay", e)
	default:
		t.Fatal(e)
	}
	switch _, e := GetBool(&run, nil); e.(type) {
	case MissingEval:
		t.Log("okay", e)
	default:
		t.Fatal(e)
	}
	switch _, e := GetNumber(&run, nil); e.(type) {
	case MissingEval:
		t.Log("okay", e)
	default:
		t.Fatal(e)
	}
	switch _, e := GetText(&run, nil); e.(type) {
	case MissingEval:
		t.Log("okay", e)
	default:
		t.Fatal(e)
	}

}
