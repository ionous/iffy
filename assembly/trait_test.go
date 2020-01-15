package assembly

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/kr/pretty"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dbutil"
	"github.com/ionous/iffy/ephemera"
)

func makeTraits(db *sql.DB, pairs []pair) (err error) {
	dbq := ephemera.NewDBQueue(db)
	r := ephemera.NewRecorder("test", dbq)
	w := NewModeler(dbq)

	for _, p := range pairs {
		var aspect, trait ephemera.Named
		if len(p.key) > 0 {
			aspect = r.Named(ephemera.NAMED_ASPECT, p.key, "key")
		}
		if len(p.value) > 0 {
			trait = r.Named(ephemera.NAMED_TRAIT, p.value, "value")
		}
		if aspect.IsValid() && trait.IsValid() {
			r.NewAspect(aspect)
			r.NewTrait(trait, aspect, 0)
		}
	}
	return DetermineAspects(w, db)
}

type expectedTrait struct {
	aspect, trait string
	rank          int
}

func matchTraits(db *sql.DB, want []expectedTrait) (err error) {
	var curr expectedTrait
	var have []expectedTrait
	if e := dbutil.QueryAll(db,
		`select aspect,trait,rank from mdl_rank order by aspect, trait, rank`,
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
	const source = memory
	if db, e := sql.Open("sqlite3", memory); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		if e := makeTraits(db,
			[]pair{
				{"A", "x"},
				{"A", "y"},
				{"B", "z"},
				{"B", "z"},
			}); e != nil {
			t.Fatal(e)
		} else if e := matchTraits(db, []expectedTrait{
			{"A", "x", 0},
			{"A", "y", 0},
			{"B", "z", 0},
		}); e != nil {
			t.Fatal("matchTraits:", e)
		} else if e := MissingAspects(db, func(a string) error {
			return errutil.New("missing aspect", a)
		}); e != nil {
			t.Fatal(e)
		} else if e := MissingTraits(db, func(t string) error {
			return errutil.New("missing trait", t)
		}); e != nil {
			t.Fatal(e)
		}
	}
}

// TestTraitConflicts
func TestTraitConflicts(t *testing.T) {
	const source = memory
	if db, e := sql.Open("sqlite3", memory); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		if e := makeTraits(db,
			[]pair{
				{"A", "x"},
				{"C", "z"},
				{"B", "x"},
			}); e == nil {
			t.Fatal("expected an error")
		} else {
			t.Log("okay:", e)
		}
	}
}
func TestTraitMissingAspect(t *testing.T) {
	const source = memory
	if db, e := sql.Open("sqlite3", memory); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		if e := makeTraits(db,
			[]pair{
				{"A", "x"},
				{"Z", ""},
			}); e != nil {
			t.Fatal(e)
		} else {
			var aspects []string
			if e := MissingAspects(db, func(a string) (err error) {
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
	const source = memory
	if db, e := sql.Open("sqlite3", memory); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		if e := makeTraits(db,
			[]pair{
				{"A", "x"},
				{"", "y"},
				{"", "z"},
			}); e != nil {
			t.Fatal(e)
		} else {
			var traits []string
			if e := MissingTraits(db, func(a string) (err error) {
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
