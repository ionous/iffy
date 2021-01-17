package decode_test

import (
	"encoding/json"
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ephemera/debug"
	"github.com/ionous/iffy/ephemera/decode"
	"github.com/kr/pretty"
)

// read simple unit test story into memory as a golang struct
func TestDecode(t *testing.T) {
	dec := decode.NewDecoder()
	// register creation functions for all the slats.
	dec.AddDefaultCallbacks(core.Slats)
	// read say story data
	var spec map[string]interface{}
	if e := json.Unmarshal([]byte(debug.SayHelloGoodbyeData), &spec); e != nil {
		t.Fatal(e)
	} else if prog, e := dec.ReadSpec(spec); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(debug.SayHelloGoodbye, prog); len(diff) > 0 {
		t.Fatal(pretty.Sprint(prog))
	} else {
		t.Log("ok. decoded story matches expected story")
	}
}
