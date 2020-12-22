package assembly

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/tables"
	"github.com/ionous/iffy/test/testdb"
	"github.com/kr/pretty"
)

// TestNounFormation to verify we can successfully assemble nouns from ephemera
func TestNounFormation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		if asm, e := newAssemblyTest(t, testdb.Memory); e != nil {
			t.Fatal(e)
		} else {
			defer asm.db.Close()
			if e := AddTestHierarchy(asm.assembler,
				"Ts", "",
			); e != nil {
				t.Fatal(e)
			} else if e := addNounEphemera(asm.rec,
				"apple", "T",
				"pear", "T",
				"toy boat", "T",
			); e != nil {
				t.Fatal(e)
			} else if e := AssembleNouns(asm.assembler); e != nil {
				t.Log(asm.dilemmas)
				t.Fatal(e)
			} else if e := matchNouns(asm.db, []modeledNoun{
				{"apple", "Ts", 0},
				{"pear", "Ts", 0},
				{"toy_boat", "Ts", 0},
				{"boat", "Ts", 1},
				{"toy", "Ts", 2},
			}); e != nil {
				t.Fatal(e)
			}
		}
	})
	t.Run("failure", func(t *testing.T) {
		if asm, e := newAssemblyTest(t, testdb.Memory); e != nil {
			t.Fatal(e)
		} else {
			defer asm.db.Close()
			if e := AddTestHierarchy(asm.assembler,
				"Ts", "",
			); e != nil {
				t.Fatal(e)
			} else if e := addNounEphemera(asm.rec,
				"bad apple", "Bs",
			); e != nil {
				t.Fatal(e)
			} else if e := AssembleNouns(asm.assembler); e == nil {
				t.Fatal("expected error")
			} else if !containsOnly(asm.dilemmas, `missing noun: "bad_apple"`) {
				t.Fatal(asm.dilemmas)
			}
		}
	})
}

func collectNouns(db *sql.DB) (ret []modeledNoun, err error) {
	var curr modeledNoun
	var nouns []modeledNoun
	if e := tables.QueryAll(db,
		`select me.name, mn.kind, me.rank
		from mdl_name me
		join mdl_noun mn
			using (noun)
		order by me.noun, me.rank, me.name`,
		func() (err error) {
			nouns = append(nouns, curr)
			return
		}, &curr.name, &curr.kind, &curr.rank); e != nil {
		err = e
	} else {
		ret = nouns
	}
	return
}

type modeledNoun struct {
	name, kind string
	rank       int
}

func addNounEphemera(rec *ephemera.Recorder, els ...string) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 2 {
		// fix: not sure its right that we breakcase the import
		// ( rather than figure break case during assembly )
		// but that's the way its done right now.
		key, value := lang.Breakcase(els[i]), lang.Breakcase(els[i+1])
		n := rec.NewName(key, tables.NAMED_NOUN, "test")
		k := rec.NewName(value, tables.NAMED_KIND, "test")
		rec.NewNoun(n, k)
	}
	return
}

// TestNounLcaSuccess to verify we can successfully determine the lowest common ancestor of nouns.
func TestNounLcaSuccess(t *testing.T) {
	if asm, e := newAssemblyTest(t, testdb.Memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		if e := AddTestHierarchy(asm.assembler,
			"Ts", "",
			"Ps", "Ts",
			"Cs", "Ps,Ts",
			"Ds", "Ps,Ts",
		); e != nil {
			t.Fatal(e)
		} else if e := addNounEphemera(asm.rec,
			"apple", "C",
			"apple", "P",
			"pear", "D",
			"pear", "T",
			"bandanna", "C",
		); e != nil {
			t.Fatal(e)
		} else if e := AssembleNouns(asm.assembler); e != nil {
			t.Fatal(e)
		} else if e := matchNouns(asm.db, []modeledNoun{
			{"apple", "Ps", 0},
			{"bandanna", "Cs", 0},
			{"pear", "Ts", 0},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

// TestNounLcaFailure to verify a mismatched noun hierarchy generates an error.
func TestNounLcaFailure(t *testing.T) {
	if asm, e := newAssemblyTest(t, testdb.Memory); e != nil {
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
		} else if e := addNounEphemera(asm.rec,
			"apple", "C",
			"apple", "D",
		); e != nil {
			t.Fatal(e)
		} else if e := AssembleNouns(asm.assembler); e == nil {
			t.Fatal("expected failure")
		} else {
			t.Log("okay:", e)
		}
	}
}

// TestNounParts to verify a single noun generates multi part names
func TestNounParts(t *testing.T) {
	if asm, e := newAssemblyTest(t, testdb.Memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		if e := AddTestHierarchy(asm.assembler,
			"Ts", "",
		); e != nil {
			t.Fatal(e)
		} else if e := addNounEphemera(asm.rec,
			"collection of words", "T",
		); e != nil {
			t.Fatal(e)
		} else if e := AssembleNouns(asm.assembler); e != nil {
			t.Fatal(e)
		} else if e := matchNouns(asm.db, []modeledNoun{
			{"collection_of_words", "Ts", 0},
			{"words", "Ts", 1},
			{"of", "Ts", 2},
			{"collection", "Ts", 3},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

func matchNouns(db *sql.DB, want []modeledNoun) (err error) {
	if got, e := collectNouns(db); e != nil {
		err = e
	} else if !reflect.DeepEqual(got, want) {
		e := errutil.New("mismatch",
			"have:", pretty.Sprint(got),
			"want:", pretty.Sprint(want))
		err = e
	}
	return
}
