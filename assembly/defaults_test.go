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

// TestDefaultFieldAssigment to verify default values can be assigned to kinds.
func TestDefaultFieldAssigment(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		if e := AddTestHierarchy(asm.assembler,
			"Ks", "",
			"Ls", "Ks",
			"Ds", "Ks",
			"Cs", "Ls,Ks",
		); e != nil {
			t.Fatal(e)
		} else if e := AddTestFields(asm.assembler,
			"Ks", "d", tables.PRIM_DIGI,
			"Ks", "t", tables.PRIM_TEXT,
			"Ks", "t2", tables.PRIM_TEXT,
			"Ls", "x", tables.PRIM_TEXT,
			"Ds", "x", tables.PRIM_TEXT,
			"Cs", "c", tables.PRIM_TEXT,
		); e != nil {
			t.Fatal(e)
		} else if e := addDefaults(asm.rec,
			"Ks", "t", "some text",
			"Ls", "t", "override text",
			"Ls", "t2", "other text",
			"Ls", "x", "x in p",
			"Ds", "x", "x in d",
			"Cs", "c", "c text",
			"Cs", "d", 123,
		); e != nil {
			t.Fatal(e)
		} else if e := AssembleDefaults(asm.assembler); e != nil {
			t.Fatal(e)
		} else if e := matchDefaults(asm.db,
			"Cs", "c", "c text",
			"Cs", "d", int64(123), // re: int64 -- default scanner uses https://golang.org/pkg/database/sql/#Scanner
			"Ds", "x", "x in d",
			"Ks", "t", "some text",
			"Ls", "t", "override text",
			"Ls", "t2", "other text",
			"Ls", "x", "x in p",
		); e != nil {
			t.Fatal(e)
		}
	}
}

// TestDefaultTraitAssignment to verify default traits can be assigned to kinds.
func TestDefaultTraitAssignment(t *testing.T) {
	if asm, e := newDefaultsTest(t, memory,
		"Ks", "x", true,
		"Ls", "y", true,
		"Ls", "z", true,
		//
		"Ns", "A", "w",
		"Ns", "B", "z",
		"Ns", "w", true,
	); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		if e := AssembleDefaults(asm.assembler); e != nil {
			t.Fatal(e)
		} else if e := matchDefaults(asm.db,
			"Ks", "A", "x",
			"Ls", "A", "y",
			"Ls", "B", "z",
			"Ns", "A", "w",
			"Ns", "B", "z",
		); e != nil {
			t.Fatal(e)
		}
	}
}

// TestDefaultDuplicates to verify that duplicate default specifications are okay
func TestDefaultDuplicates(t *testing.T) {
	if asm, e := newDefaultsTest(t, memory,
		"Ks", "t", "text",
		"Ks", "t", "text",
		"Ls", "t", "text",
		//
		"Ks", "d", 123,
		"Ks", "d", 123,
		"Ls", "d", 123,
		//
		"Ks", "A", "y",
		"Ks", "y", true,
		"Ls", "x", true,
		"Ls", "A", "x",
	); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		if e := AssembleDefaults(asm.assembler); e != nil {
			t.Fatal(e)
		}
	}
}

// TestDefaultConflict to verify that conflicting values for the same default are not okay
func TestDefaultConflict(t *testing.T) {
	testConflict := func(t *testing.T, vals ...interface{}) (err error) {
		if asm, e := newDefaultsTest(t, memory, vals...); e != nil {
			t.Fatal(e)
		} else {
			defer asm.db.Close()
			if e := AssembleDefaults(asm.assembler); e == nil {
				err = errutil.New("expected error")
			} else {
				t.Log("okay:", e)
			}
		}
		return
	}
	if e := testConflict(t,
		"Ks", "t", "a",
		"Ks", "t", "b",
	); e != nil {
		t.Fatal(e)
	} else if e := testConflict(t,
		"Ks", "d", 1,
		"Ks", "d", 2,
	); e != nil {
		t.Fatal(e)
	}
	if e := testConflict(t,
		"Ks", "A", "x",
		"Ks", "A", "y",
	); e != nil {
		t.Fatal(e)
	} else if e := testConflict(t,
		"Ks", "x", true,
		"Ks", "y", true,
	); e != nil {
		t.Fatal(e)
	} else if e := testConflict(t,
		"Ks", "A", "x",
		"Ks", "y", true,
	); e != nil {
		t.Fatal(e)
	}
}

// TestDefaultBadValue to verify that modeling requires appropriate values for defaults based on type
func TestDefaultBadValue(t *testing.T) {
	//- for now, we only allow text and number [ text and number ]
	// - later we could add ambiguity for conversion [ 4 -> "4" ]
	testInvalid := func(t *testing.T, vals ...interface{}) (err error) {
		if asm, e := newDefaultsTest(t, memory, vals...); e != nil {
			err = e
		} else {
			defer asm.db.Close()
			if e := AssembleDefaults(asm.assembler); e == nil {
				err = errutil.New("expected error")
			} else {
				t.Log("okay:", e)
			}
		}
		return
	}

	if e := testInvalid(t,
		"Ks", "t", 1.2,
	); e != nil {
		t.Fatal(e)
	} else if e := testInvalid(t,
		"Ks", "d", "1.2",
	); e != nil {
		t.Fatal(e)
	}
	// try to set trait like values
	if e := testInvalid(t,
		"Ks", "t", false,
	); e != nil {
		t.Fatal(e)
	}

	/*
	   fix? somehow? bools in sqlite are stored as int64;
	   could switch to text ( "true", "false" ) perhaps and add some check/query
	   during determination
	   if e := testInvalid(t,
	       "Ks", "d", true,
	   ); e != nil {
	       t.Fatal(e)
	   }
	*/

	/* fix? aspects are set by matching traits
	1.2 is not a trait, so it's skipped.
	this might get handled by a "missing" check,
	or possibly by changing the determination query.

	if e := testInvalid(t,
		"Ks", "A", 1.2,
	); e != nil {
		t.Fatal(e)
	}
	*/
}

// match generated model defaults
func matchDefaults(db *sql.DB, want ...interface{}) (err error) {
	var a, b, c interface{}
	var have []interface{}
	if e := tables.QueryAll(db,
		`select kind, field, value 
			from mdl_default
			order by kind, field, value`,
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

// write ephemera describing some initial values
func addDefaults(rec *ephemera.Recorder, els ...interface{}) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 3 {
		tgt, prop, val := els[i], els[i+1], els[i+2]
		namedKind := rec.NewName(tgt.(string), tables.NAMED_KINDS, "test")
		namedField := rec.NewName(prop.(string), tables.NAMED_PROPERTY, "test")
		rec.NewDefault(namedKind, namedField, val)
	}
	return
}

func newDefaultsTest(t *testing.T, path string, defaults ...interface{}) (ret *assemblyTest, err error) {
	ret = &assemblyTest{T: t}
	if asm, e := newAssemblyTest(t, path); e != nil {
		err = e
	} else {
		if e := AddTestHierarchy(asm.assembler,
			"Ks", "",
			"Ls", "Ks",
			"Ns", "Ks",
		); e != nil {
			err = e
		} else if e := AddTestFields(asm.assembler,
			"Ks", "d", tables.PRIM_DIGI,
			"Ks", "t", tables.PRIM_TEXT,
			"Ks", "A", tables.PRIM_ASPECT,
			"Ls", "B", tables.PRIM_ASPECT,
			"Ns", "B", tables.PRIM_ASPECT,
		); e != nil {
			err = e
		} else if e := AddTestTraits(asm.assembler,
			"A", "w",
			"A", "x",
			"A", "y",
			"B", "z",
		); e != nil {
			err = e
		} else if e := addDefaults(asm.rec, defaults...); e != nil {
			err = e
		}
		if err != nil {
			asm.db.Close()
		} else {
			ret = asm
		}
	}
	return
}
