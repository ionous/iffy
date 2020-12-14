package story

import (
	"testing"
)

// read simple unit test story into sqlite, extract as a golang struct
// fix.... tests require not just output but rules now.
func xTestProcessProg(t *testing.T) {
	// k, db := newImporter(t, testdb.Memory)
	// defer db.Close()
	// if e := imp_test_output(k, k.NewName("test", t.Name(), ""), debug.SayHelloGoodbyeData); e != nil {
	// 	t.Fatal(e)
	// } else {
	// 	var testName string
	// 	var progid int
	// 	var expect string
	// 	if e := db.QueryRow("select * from eph_check").Scan(&testName, &progid, &expect); e != nil {
	// 		t.Fatal(e)
	// 	} else {
	// 		t.Log(testName, progid, expect)
	// 		//
	// 		var id int
	// 		var typeName string
	// 		var prog []byte
	// 		if e := db.QueryRow("select * from eph_prog").Scan(&id, &typeName, &prog); e != nil {
	// 			t.Fatal(e)
	// 		} else {
	// 			var res check.Testing
	// 			t.Log(id, typeName)
	// 			dec := gob.NewDecoder(bytes.NewBuffer(prog))
	// 			if e := dec.Decode(&res); e != nil {
	// 				t.Fatal(e)
	// 			} else if diff := pretty.Diff(&debug.SayHelloGoodbye, res); len(diff) > 0 {
	// 				t.Fatal(diff)
	// 			}
	// 		}
	// 	}
	// }
}
