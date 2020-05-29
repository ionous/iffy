package assembly

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/ionous/iffy/tables"
)

// TestMissingKinds to verify the kinds mentioned in parent-child ephemera exist.
func TestMissingKinds(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		// kind, ancestor
		pairs := []string{
			"P", "T",
			"Q", "T",
			"P", "R",
		}
		for i := 0; i < len(pairs); i += 2 {
			kid := asm.rec.NewName(pairs[i], tables.NAMED_KINDS, strconv.Itoa(i))
			parent := asm.rec.NewName(pairs[i+1], tables.NAMED_KINDS, strconv.Itoa(i+1))
			asm.rec.NewKind(kid, parent)
		}
		if e := AssembleAncestry(asm.assembler, "T"); e == nil {
			t.Fatal("expected error")
		} else if !containsOnly(asm.dilemmas, `missing kind: "R"`) {
			t.Fatal(asm.dilemmas)
		} else {
			t.Log("ok:", e)
		}
	}
}

// TestMissingAspects detects fields labeled as aspects which are missing from the aspects ephemera.
func TestMissingAspects(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		parent := asm.rec.NewName("K", tables.NAMED_KINDS, "container")
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
			asm.rec.NewPrimitive(parent, a, tables.PRIM_ASPECT)
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
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		if e := AddTestHierarchy(asm.assembler, []TargetField{
			{"T", ""},
		}); e != nil {
			t.Fatal(e)
		} else if e := writeMissing(asm.rec, []string{
			"z",
		}); e != nil {
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
// 	if asm, e := newAssemblyTest(t, memory); e != nil {
// 		t.Fatal(e)
// 	} else {
// 		defer asm.db.Close()
// 		//
// 		if e := AddTestHierarchy(asm.assembler, []TargetField{
// 			{"T", ""},
// 			{"P", "T"},
// 			{"C", "P,T"},
// 		}); e != nil {
// 			t.Fatal(e)
// 		} else if e := AddTestFields(asm.assembler, []TargetValue{
// 			{"T", "d", tables.PRIM_DIGI},
// 			{"T", "t", tables.PRIM_TEXT},
// 			{"T", "t2", tables.PRIM_TEXT},
// 			{"P", "p", tables.PRIM_TEXT},
// 			{"C", "c", tables.PRIM_TEXT},
// 		}); e != nil {
// 			t.Fatal(e)
// 		} else if e := addDefaults(asm.rec, []triplet{
// 			{"T", "t", "some text"},
// 			{"P", "t2", "other text"},
// 			{"C", "c", "c text"},
// 			{"P", "c", "invalid"}, // this pair doesnt exist
// 			{"T", "p", "invalid"}, // this pair doesnt exist
// 			{"C", "d", 123},
// 		}); e != nil {
// 			t.Fatal(e)
// 		} else {
// 			var got []pair
// 			if e := MissingDefaults(asm.db, func(k, f string) (err error) {
// 				got = append(got, pair{k, f})
// 				return
// 			}); e != nil {
// 				t.Fatal(e)
// 			} else if !reflect.DeepEqual(got, []pair{
// 				{"P", "c"},
// 				{"T", "p"},
// 			}) {
// 				t.Fatal("mismatched", got)
// 			}
// 		}
// 	}
// }
