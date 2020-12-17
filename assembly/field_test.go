package assembly

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/kr/pretty"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
	"github.com/ionous/iffy/test/testdb"
)

// write some primitives
func writeFields(rec *ephemera.Recorder, els ...string) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 3 {
		kind, field, fieldType := els[i], els[i+1], els[i+2]
		kn := rec.NewName(kind, tables.NAMED_KINDS, "test")
		fn := rec.NewName(field, tables.NAMED_FIELD, "test")
		rec.NewField(kn, fn, fieldType, "")
	}
	return
}

// name some fields that arent otherwise referenced
func writeMissing(rec *ephemera.Recorder, missing ...string) (err error) {
	for _, m := range missing {
		rec.NewName(m, tables.NAMED_FIELD, "test")
	}
	return
}

func matchProperties(db *sql.DB, want ...string) (err error) {
	var kind, field, fieldType string
	var have []string
	if e := tables.QueryAll(db,
		`select kind, field, type 
		from mdl_field 
		order by kind, field, type`,
		func() (err error) {
			have = append(have, kind, field, fieldType)
			return
		}, &kind, &field, &fieldType); e != nil {
		err = e
	} else if !reflect.DeepEqual(have, want) {
		err = errutil.New("mismatch", "have:", pretty.Sprint(have), "want:", pretty.Sprint(want))
	}
	return
}

func TestFields(t *testing.T) {
	if asm, e := newAssemblyTest(t, testdb.Memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		if e := AddTestHierarchy(asm.assembler,
			"Ts", "",
			"Ps", "Ts",
			"Qs", "Ts",
		); e != nil {
			t.Fatal(e)
		} else if e := writeFields(asm.rec,
			"Ps", "a", tables.PRIM_TEXT,
			"Qs", "b", tables.PRIM_TEXT,
			"Ts", "c", tables.PRIM_TEXT,
		); e != nil {
			t.Fatal(e)
		} else if e := AssembleFields(asm.assembler); e != nil {
			t.Fatal(e)
		} else if e := matchProperties(asm.db,
			"Ps", "a", tables.PRIM_TEXT,
			"Qs", "b", tables.PRIM_TEXT,
			"Ts", "c", tables.PRIM_TEXT,
		); e != nil {
			t.Fatal(e)
		}
	}
}

func TestFieldLca(t *testing.T) {
	if asm, e := newAssemblyTest(t, testdb.Memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		if e := AddTestHierarchy(asm.assembler,
			"Ts", "",
			"Ps", "Ts",
			"Qs", "Ts",
		); e != nil {
			t.Fatal(e)
		} else if e := writeFields(asm.rec,
			"Ps", "a", tables.PRIM_TEXT,
			"Qs", "a", tables.PRIM_TEXT,
		); e != nil {
			t.Fatal(e)
		} else if e := AssembleFields(asm.assembler); e != nil {
			t.Fatal(e)
		} else if e := matchProperties(asm.db,
			"Ts", "a", tables.PRIM_TEXT,
		); e != nil {
			t.Fatal(e)
		}
	}
}

// TestFieldTypeMismatch verifies that ephemera with conflicting primitive types generates an error
// ex. T.a:text, T.a:number
func TestFieldTypeMismatch(t *testing.T) {
	if asm, e := newAssemblyTest(t, testdb.Memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		if e := AddTestHierarchy(asm.assembler,
			"Ts", "",
		); e != nil {
			t.Fatal(e)
		} else if e := writeFields(asm.rec,
			"Ts", "a", tables.PRIM_TEXT,
			"Ts", "a", tables.PRIM_DIGI,
		); e != nil {
			t.Fatal(e)
		} else if e := AssembleFields(asm.assembler); e != nil {
			t.Log("okay:", e)
		} else {
			t.Fatal("expected error")
		}
	}
}
