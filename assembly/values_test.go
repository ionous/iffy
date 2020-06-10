package assembly

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
	"github.com/kr/pretty"
)

// TestInitialFieldAssignment to verify initial values for fields can be assigned to instances.
func TestInitialFieldAssignment(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		if e := AddTestHierarchy(asm.assembler,
			"Ks", "",
			"Ls", "Ks",
			"Ms", "Ls,Ks",
		); e != nil {
			t.Fatal(e)
		} else if e := AddTestFields(asm.assembler,
			"Ks", "t", tables.PRIM_TEXT,
			"Ls", "d", tables.PRIM_DIGI,
		); e != nil {
			t.Fatal(e)
		} else if e := AddTestNouns(asm.assembler,
			"apple", "Ks",
			"pear", "Ls",
			"toy boat", "Ms",
			"boat", "Ms",
		); e != nil {
			t.Fatal(e)
		} else if e := addValues(asm.rec,
			"apple", "t", "some text",
			"pear", "d", 123,
			"toy", "d", 321,
			"boat", "t", "more text",
		); e != nil {
			t.Fatal(e)
		} else if e := assembleInitialFields(asm.assembler); e != nil {
			t.Fatal(e)
		} else if e := matchValues(asm.db,
			"apple", "t", "some text",
			"boat", "t", "more text",
			"pear", "d", int64(123), // int64, re: go's default scanner.
			"toy boat", "d", int64(321),
		); e != nil {
			t.Fatal(e)
		}
	}
}

// TestInitialTraitAssignments to verify default traits can be assigned to kinds.
func TestInitialTraitAssignment(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		if e := AddTestHierarchy(asm.assembler,
			"Ks", "",
			"Ls", "Ks",
			"Ms", "Ls,Ks",
		); e != nil {
			t.Fatal(e)
		} else if e := AddTestFields(asm.assembler,
			"Ks", "A", tables.PRIM_ASPECT,
			"Ls", "B", tables.PRIM_ASPECT,
		); e != nil {
			t.Fatal(e)
		} else if e := AddTestTraits(asm.assembler,
			"A", "w",
			"A", "x",
			"A", "y",
			"B", "z",
		); e != nil {
			t.Fatal(e)
		} else if e := AddTestNouns(asm.assembler,
			"apple", "Ks",
			"pear", "Ls",
			"toy boat", "Ms",
			"boat", "Ms",
		); e != nil {
			t.Fatal(e)
		} else if e := addValues(asm.rec,
			"apple", "A", "y",
			"pear", "x", true,
			"toy", "w", true,
			"boat", "z", true,
		); e != nil {
			t.Fatal(e)
		} else if e := AssembleValues(asm.assembler); e != nil {
			t.Fatal(e)
		} else if e := matchValues(asm.db,
			"apple", "A", "y",
			"boat", "B", "z",
			"pear", "A", "x",
			"toy boat", "A", "w",
		); e != nil {
			t.Fatal(e)
		}
	}
}

// match generated model defaults
func matchValues(db *sql.DB, want ...interface{}) (err error) {
	var a, b, c interface{}
	var have []interface{}
	if e := tables.QueryAll(db,
		`select noun, field, value 
			from mdl_start
			order by noun, field, value`,
		func() (err error) {
			have = append(have, a, b, c)
			return
		},
		&a, &b, &c); e != nil {
		err = e
	} else if !reflect.DeepEqual(have, want) {
		err = errutil.New("mismatch",
			"have:", pretty.Sprint(have),
			"want:", pretty.Sprint(want))
	}
	return
}

// eph_value: fake noun, prop, value
// prop: k, f, v
func addValues(rec *ephemera.Recorder, els ...interface{}) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 3 {
		tgt, field, value := els[i], els[i+1], els[i+2]
		noun := rec.NewName(tgt.(string), tables.NAMED_NOUN, "test")
		prop := rec.NewName(field.(string), tables.NAMED_PROPERTY, "test")
		rec.NewValue(noun, prop, value)
	}
	return
}
