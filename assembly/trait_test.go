package assembly

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/kr/pretty"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dbutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
)

func addTraits(rec *ephemera.Recorder, pairs []pair) (err error) {
	for _, p := range pairs {
		var aspect, trait ephemera.Named
		if len(p.key) > 0 {
			aspect = rec.Named(tables.NAMED_ASPECT, p.key, "key")
		}
		if len(p.value) > 0 {
			trait = rec.Named(tables.NAMED_TRAIT, p.value, "value")
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
	if e := dbutil.QueryAll(db,
		`select aspect, trait, rank from mdl_aspect order by aspect, trait, rank`,
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
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		//
		if e := addTraits(t.rec,
			[]pair{
				{"A", "x"},
				{"A", "y"},
				{"B", "z"},
				{"B", "z"},
			}); e != nil {
			t.Fatal(e)
		} else if e := DetermineAspects(t.modeler, t.db); e != nil {
			t.Fatal(e)
		} else if e := matchTraits(t.db, []expectedTrait{
			{"A", "x", 0},
			{"A", "y", 0},
			{"B", "z", 0},
		}); e != nil {
			t.Fatal("matchTraits:", e)
		} else if e := MissingAspects(t.db, func(a string) error {
			return errutil.New("missing aspect", a)
		}); e != nil {
			t.Fatal(e)
		} else if e := MissingTraits(t.db, func(t string) error {
			return errutil.New("missing trait", t)
		}); e != nil {
			t.Fatal(e)
		}
	}
}

// TestTraitConflicts
func TestTraitConflicts(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		//
		if e := addTraits(t.rec, []pair{
			{"A", "x"},
			{"C", "z"},
			{"B", "x"},
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineAspects(t.modeler, t.db); e == nil {
			t.Fatal("expected an error")
		} else {
			t.Log("okay:", e)
		}
	}
}
func TestTraitMissingAspect(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		//
		if e := addTraits(t.rec, []pair{
			{"A", "x"},
			{"Z", ""},
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineAspects(t.modeler, t.db); e != nil {
			t.Fatal(e)
		} else {
			var aspects []string
			if e := MissingAspects(t.db, func(a string) (err error) {
				aspects = append(aspects, a)
				return
			}); e != nil {
				t.Fatal(e)
			} else if !reflect.DeepEqual(aspects, []string{"Z"}) {
				t.Fatal("mismatch", aspects)
			} else {
				t.Log("okay, missing", aspects)
			}
		}
	}
}

func TestTraitMissingTraits(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		//
		if e := addTraits(t.rec, []pair{
			{"A", "x"},
			{"", "y"},
			{"", "z"},
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineAspects(t.modeler, t.db); e != nil {
			t.Fatal(e)
		} else {
			var traits []string
			if e := MissingTraits(t.db, func(a string) (err error) {
				traits = append(traits, a)
				return
			}); e != nil {
				t.Fatal(e)
			} else if !reflect.DeepEqual(traits, []string{"y", "z"}) {
				t.Fatal("mismatch", traits)
			} else {
				t.Log("okay, missing:", traits)
			}
		}
	}
}
