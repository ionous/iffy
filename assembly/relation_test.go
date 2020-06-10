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

// TestVerbMismatches verifies that we can collapse multiple relation-verb pairs so long as the verb-stem pair match
// while ensuring the same stem cannot be used in multiple relations.
func TestVerbMismatches(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal("error creating test", e)
	} else {
		defer asm.db.Close()
		m := asm.assembler
		//
		if e := m.WriteVerb("R", "contains"); e != nil {
			t.Fatal(e)
		} else if e := m.WriteVerb("R", "containing"); e != nil {
			t.Fatal(e)
		} else if e := m.WriteVerb("Q", "supporting"); e != nil {
			t.Fatal(e)
		} else if e := m.WriteVerb("Q", "supports"); e != nil {
			t.Fatal(e)
		} else if e := m.WriteVerb("R", "supports"); e == nil {
			t.Log("expected error")
		} else {
			t.Log("okay:", e)
		}
	}
}

// write some relation ephemera to the database:
// relation, kind1, kind2, cardinality
// ( from which assembly will determine relations )
func addRelations(rec *ephemera.Recorder, els ...string) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 4 {
		er, ek, eq, ec := els[i], els[i+1]+"s", els[i+2]+"s", els[i+3]
		r := rec.NewName(er, tables.NAMED_RELATION, "test")
		k := rec.NewName(ek, tables.NAMED_KINDS, "test")
		q := rec.NewName(eq, tables.NAMED_KINDS, "test")
		c := ec
		rec.NewRelation(r, k, q, c)
	}
	return
}

// relation, kind1, kind2, cardinality
func matchRelations(db *sql.DB, want ...string) (err error) {
	var have []string
	var a, b, c, d string
	if e := tables.QueryAll(db,
		`select relation, kind, otherKind, cardinality
			from mdl_rel 
		order by relation, kind, otherKind, cardinality`,
		func() (err error) {
			have = append(have, a, b, c, d)
			return
		}, &a, &b, &c, &d); e != nil {
		err = e
	} else if !reflect.DeepEqual(have, want) {
		err = errutil.New("mismatch", "have:", pretty.Sprint(have), "want:", pretty.Sprint(want))
	}
	return
}

// TestRelationCreation to verify it's possible to build relations
func TestRelationCreation(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
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
		} else if e := addRelations(
			asm.rec,
			"R", "P", "Q", tables.ONE_TO_MANY,
			"G", "P", "Q", tables.MANY_TO_ONE,
			"G", "P", "Q", tables.MANY_TO_ONE,
			"H", "P", "P", tables.ONE_TO_MANY,
		); e != nil {
			t.Fatal(e)
		} else if e := AssembleRelations(asm.assembler); e != nil {
			t.Fatal(e)
		} else if e := matchRelations(asm.db,
			"G", "Ps", "Qs", tables.MANY_TO_ONE,
			"H", "Ps", "Ps", tables.ONE_TO_MANY,
			"R", "Ps", "Qs", tables.ONE_TO_MANY,
		); e != nil {
			t.Fatal(e)
		}
	}
}

// TestRelationCardinality detects conflicting cardinalities
func TestRelationCardinality(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		if e := AddTestHierarchy(asm.assembler,
			"Ps", "",
		); e != nil {
			t.Fatal(e)
		} else if e := addRelations(asm.rec,
			"R", "P", "P", tables.ONE_TO_MANY,
			"R", "P", "P", tables.MANY_TO_ONE,
		); e != nil {
			t.Fatal(e)
		} else if e := AssembleRelations(asm.assembler); e == nil {
			t.Fatal("expected error")
		} else {
			t.Log("okay:", e)
		}
	}
}

// TestRelationLca
func TestRelationLcaSuccess(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		if e := AddTestHierarchy(asm.assembler,
			"Ts", "",
			"Ps", "Ts",
			"Cs", "Ps,Ts",
			"Ds", "Ps,Ts",
		); e != nil {
			t.Fatal(e)
		} else if e := addRelations(asm.rec,
			"R", "P", "T", tables.ONE_TO_MANY,
			"R", "D", "T", tables.ONE_TO_MANY,
			"R", "C", "T", tables.ONE_TO_MANY,
		); e != nil {
			t.Fatal(e)
		} else if e := AssembleRelations(asm.assembler); e != nil {
			t.Fatal(e)
		} else if e := matchRelations(asm.db,
			"R", "Ps", "Ts", tables.ONE_TO_MANY,
		); e != nil {
			t.Fatal(e)
		}
	}
}

// TestRelationLcaFailure to verify a mismatched relation hierarchy generates an error.
func TestRelationLcaFailure(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		if e := AddTestHierarchy(asm.assembler,
			"Ts", "",
			"Ps", "Ts",
			"Cs", "Ps,Ts",
			"Ds", "Ps,Ts",
		); e != nil {
			t.Fatal(e)
		} else if e := addRelations(asm.rec,
			"R", "D", "T", tables.ONE_TO_MANY,
			"R", "C", "T", tables.ONE_TO_MANY,
		); e != nil {
			t.Fatal(e)
		} else if e := AssembleRelations(asm.assembler); e == nil {
			t.Fatal("expected error")
		} else {
			t.Log("okay:", e)
		}
	}
}
