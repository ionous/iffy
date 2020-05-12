package story

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/ionous/iffy/dl/check"
	"github.com/ionous/iffy/ephemera/debug"
	"github.com/kr/pretty"
)

// read simple unit test story into sqlite, extract as a golang struct
func TestProcessProg(t *testing.T) {
	k, db := newTestDecoder(t)
	defer db.Close()
	if e := imp_test_output(k, k.NewName("test", t.Name(), ""), debug.SayStory); e != nil {
		t.Fatal(e)
	} else {
		var testName string
		var progid int
		var expect string
		if e := db.QueryRow("select * from eph_check").Scan(&testName, &progid, &expect); e != nil {
			t.Fatal(e)
		} else {
			t.Log(testName, progid, expect)
			//
			var id int
			var typeName string
			var prog []byte
			if e := db.QueryRow("select * from eph_prog").Scan(&id, &typeName, &prog); e != nil {
				t.Fatal(e)
			} else {
				var res check.TestOutput
				t.Log(id, typeName)
				dec := gob.NewDecoder(bytes.NewBuffer(prog))
				if e := dec.Decode(&res); e != nil {
					t.Fatal(e)
				} else if diff := pretty.Diff(debug.SayTest, res); len(diff) > 0 {
					t.Fatal(diff)
				}
			}
		}
	}
}
