package assembly

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/kr/pretty"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
)

// write some primitives
func writeFields(rec *ephemera.Recorder, kfps []kfp) (err error) {
	for _, p := range kfps {
		kind := rec.NewName(p.kind, tables.NAMED_KINDS, "test")
		field := rec.NewName(p.field, tables.NAMED_FIELD, "test")
		rec.NewPrimitive(kind, field, p.fieldType)
	}
	return
}

// name some fields that arent otherwise referenced
func writeMissing(rec *ephemera.Recorder, missing []string) (err error) {
	for _, m := range missing {
		rec.NewName(m, tables.NAMED_FIELD, "test")
	}
	return
}

func matchProperties(db *sql.DB, want []kfp) (err error) {
	var curr kfp
	var have []kfp
	if e := tables.QueryAll(db,
		`select kind, field, type 
		from mdl_field 
		order by kind, field, type`,
		func() (err error) {
			have = append(have, curr)
			return
		}, &curr.kind, &curr.field, &curr.fieldType); e != nil {
		err = e
	} else if !reflect.DeepEqual(have, want) {
		err = errutil.New("mismatch", "have:", pretty.Sprint(have), "want:", pretty.Sprint(want))
	}
	return
}

func TestFields(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		if e := AddTestHierarchy(asm.assembler, []TargetField{
			{"T", ""},
			{"P", "T"},
			{"Q", "T"},
		}); e != nil {
			t.Fatal(e)
		} else if e := writeFields(asm.rec, []kfp{
			{"P", "a", tables.PRIM_TEXT},
			{"Q", "b", tables.PRIM_TEXT},
			{"T", "c", tables.PRIM_TEXT},
		}); e != nil {
			t.Fatal(e)
		} else if e := AssembleFields(asm.assembler); e != nil {
			t.Fatal(e)
		} else if e := matchProperties(asm.db, []kfp{
			{"P", "a", tables.PRIM_TEXT},
			{"Q", "b", tables.PRIM_TEXT},
			{"T", "c", tables.PRIM_TEXT},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

func TestFieldLca(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		if e := AddTestHierarchy(asm.assembler, []TargetField{
			{"T", ""},
			{"P", "T"},
			{"Q", "T"},
		}); e != nil {
			t.Fatal(e)
		} else if e := writeFields(asm.rec, []kfp{
			{"P", "a", tables.PRIM_TEXT},
			{"Q", "a", tables.PRIM_TEXT},
		}); e != nil {
			t.Fatal(e)
		} else if e := AssembleFields(asm.assembler); e != nil {
			t.Fatal(e)

		} else if e := matchProperties(asm.db, []kfp{
			{"T", "a", tables.PRIM_TEXT},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

// TestFieldTypeMismatch verifies that ephemera with conflicting primitive types generates an error
// ex. T.a:text, T.a:digi
func TestFieldTypeMismatch(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		if e := AddTestHierarchy(asm.assembler, []TargetField{
			{"T", ""},
		}); e != nil {
			t.Fatal(e)
		} else if e := writeFields(asm.rec, []kfp{
			{"T", "a", tables.PRIM_TEXT},
			{"T", "a", tables.PRIM_DIGI},
		}); e != nil {
			t.Fatal(e)
		} else if e := AssembleFields(asm.assembler); e != nil {
			t.Log("okay:", e)
		} else {
			t.Fatal("expected error")
		}
	}
}
