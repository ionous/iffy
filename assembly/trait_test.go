package assembly

import (
	"database/sql"
	"reflect"
	"strings"
	"testing"

	"github.com/kr/pretty"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
)

func addTraits(rec *ephemera.Recorder, pairs []string) (err error) {
	els := pairs
	for i, cnt := 0, len(els); i < cnt; i += 2 {
		key, value := els[i], els[i+1]
		var aspect, trait ephemera.Named
		if len(key) > 0 {
			aspect = rec.NewName(key, tables.NAMED_ASPECT, "key")
		}
		if len(value) > 0 {
			trait = rec.NewName(value, tables.NAMED_TRAIT, "value")
		}
		if aspect.IsValid() && trait.IsValid() {
			rec.NewAspect(aspect)
			rec.NewTrait(trait, aspect, 0)
		}
	}
	return
}

type expectedTrait struct {
	aspect, trait string
	rank          int
}

func matchTraits(db *sql.DB, want []expectedTrait) (err error) {
	var curr expectedTrait
	var have []expectedTrait
	if e := tables.QueryAll(db,
		`select aspect, trait, rank from mdl_aspect order by aspect, rank`,
		func() (err error) {
			have = append(have, curr)
			return
		}, &curr.aspect, &curr.trait, &curr.rank); e != nil {
		err = e
	} else if !reflect.DeepEqual(have, want) {
		err = errutil.New("mismatch", "have:", pretty.Sprint(have), "want:", pretty.Sprint(want))
	}
	return
}

// TestTraits to verify that aspects/traits in ephemera can become part of the model.
func TestTraits(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		if e := addTraits(asm.rec, []string{
			"A", "x",
			"A", "y",
			"B", "z",
			"B", "z",
		}); e != nil {
			t.Fatal(e)
		} else if e := AssembleAspects(asm.assembler); e != nil {
			t.Fatal(e)
		} else if e := matchTraits(asm.db, []expectedTrait{
			{"A", "x", 0},
			{"A", "y", 1},
			{"B", "z", 0},
		}); e != nil {
			t.Fatal("matchTraits:", e)
		}
	}
}

// TestTraitConflicts
func TestTraitConflicts(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		if e := addTraits(asm.rec, []string{
			"A", "x",
			"C", "z",
			"B", "x",
		}); e != nil {
			t.Fatal(e)
		} else if e := AssembleAspects(asm.assembler); e == nil {
			t.Fatal("expected an error")
		} else {
			t.Log("okay:", e)
		}
	}
}

func TestTraitMissingAspect(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		if e := addTraits(asm.rec, []string{
			"A", "x",
			"Z", "",
		}); e != nil {
			t.Fatal(e)
		} else if e := AssembleAspects(asm.assembler); e == nil {
			t.Fatal("expected error")
		} else if asm.dilemmas.Len() != 1 ||
			!strings.Contains((*asm.dilemmas)[0].Msg, `missing aspect: "Z"`) {
			t.Fatal(asm.dilemmas)
		} else {
			t.Log("ok:", e)

		}
	}
}

func TestTraitMissingTraits(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		if e := addTraits(asm.rec, []string{
			"A", "x",
			"", "y",
			"", "z",
		}); e != nil {
			t.Fatal(e)
		} else if e := AssembleAspects(asm.assembler); e == nil {
			t.Fatal("expected error")
		} else if !containsOnly(asm.dilemmas,
			`missing trait: "y"`,
			`missing trait: "z"`) {
			t.Fatal(asm.dilemmas)
		} else {
			t.Log("ok:", e)
		}
	}
}

func containsOnly(ds *Dilemmas, msg ...string) bool {
	return ds.Len() == len(msg) && containsMessages(ds, msg...)
}

func containsMessages(ds *Dilemmas, msg ...string) (ret bool) {
	for _, d := range *ds {
		foundAt := -1
		for i, str := range msg {
			if strings.Contains(d.Msg, str) {
				foundAt = i
				break
			}
		}
		if foundAt >= 0 {
			if end := len(msg) - 1; end == 0 {
				ret = true
				break
			} else {
				// cut w/o preserving order
				msg[foundAt] = msg[end]
				msg = msg[:end]
			}
		}

	}
	return
}
