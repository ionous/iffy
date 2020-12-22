package assembly

import (
	"reflect"
	"testing"

	"github.com/ionous/iffy/tables"
	"github.com/ionous/iffy/test/testdb"
)

// TestMissingKinds to verify the kinds mentioned in parent-child ephemera exist.
func TestMissingKinds(t *testing.T) {
	if asm, e := newAssemblyTest(t, testdb.Memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		addKinds(asm, // kid, ancestor
			"T", "",
			"P", "T",
			"Q", "T",
			"P", "R",
		)
		if e := AssembleAncestry(asm.assembler, "Ts"); e == nil {
			t.Fatal("expected error")
		} else if !containsOnly(asm.dilemmas, `missing singular_kind: "R"`) {
			if d := asm.dilemmas; d.Len() == 0 {
				t.Fatal(e)
			} else {
				t.Fatal(d)
			}
		} else {
			t.Log("ok:", e)
		}
	}
}

// TestMissingAspects detects fields labeled as aspects which are missing from the aspects ephemera.
func TestMissingAspects(t *testing.T) {
	if asm, e := newAssemblyTest(t, testdb.Memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		parent := asm.rec.NewName("Ks", tables.NAMED_KINDS, "container")
		for i, aspect := range []string{
			// known, unknown
			"A", "F",
			"C", "D",
			"E", "B",
		} {
			a := asm.rec.NewName(aspect, tables.NAMED_ASPECT, "test")
			if known := i&1 == 0; known {
				asm.rec.NewAspect(a)
			}
			asm.rec.NewField(parent, a, tables.PRIM_ASPECT, "")
		}
		expected := []string{"B", "D", "F"}
		if missing, e := undeclaredAspects(asm.db); e != nil {
			t.Fatal(e)
		} else if matches := reflect.DeepEqual(missing, expected); !matches {
			t.Fatal("want:", expected, "have:", missing)
		} else {
			t.Log("okay")
		}
	}
}

func TestMissingField(t *testing.T) {
	if asm, e := newAssemblyTest(t, testdb.Memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		if e := AddTestHierarchy(asm.assembler,
			"Ts", "",
		); e != nil {
			t.Fatal(e)
		} else if e := writeMissing(asm.rec,
			"z",
		); e != nil {
			t.Fatal(e)
		} else if e := AssembleFields(asm.assembler); e == nil {
			t.Fatal("expected error")
		} else if !containsOnly(asm.dilemmas, `missing field: "z"`) {
			t.Fatal(asm.dilemmas)
		} else {
			t.Log("ok:", e)
		}
	}
}

// xTestMissingUnknownField missing properties ( kind, field pair doesn't exist in model )
// func xTestMissingUnknownField(t *testing.T) {
// 	if asm, e := newAssemblyTest(t, testdb.Memory); e != nil {
// 		t.Fatal(e)
// 	} else {
// 		defer asm.db.Close()
// 		//
// 		if e := AddTestHierarchy(asm.assembler,
// 			"Ts", "",
// 			"Ps", "Ts",
// 			"Cs", "Ps,Ts",
// 		); e != nil {
// 			t.Fatal(e)
// 		} else if e := AddTestFields(asm.assembler,
// 			"Ts", "d", tables.PRIM_DIGI,"",
// 			"Ts", "t", tables.PRIM_TEXT,"",
// 			"Ts", "t2", tables.PRIM_TEXT,"",
// 			"Ps", "p", tables.PRIM_TEXT,"",
// 			"Cs", "c", tables.PRIM_TEXT,"",
// 		); e != nil {
// 			t.Fatal(e)
// 		} else if e := addDefaults(asm.rec,
// 			"T", "t", "some text",
// 			"P", "t2", "other text",
// 			"C", "c", "c text",
// 			"P", "c", "invalid", // this pair doesnt exist
// 			"T", "p", "invalid", // this pair doesnt exist
// 			"C", "d", 123,
// 		}); e != nil {
// 			t.Fatal(e)
// 		} else {
// 			var got []interface{}
// 			if e := MissingDefaults(asm.db, func(k, f string) (err error) {
// 				got = append(got, k, f)
// 				return
// 			); e != nil {
// 				t.Fatal(e)
// 			} else if !reflect.DeepEqual(got,
// 				"P", "c",
// 				"T", "p",
// 			) {
// 				t.Fatal("mismatched", got)
// 			}
// 		}
// 	}
// }
