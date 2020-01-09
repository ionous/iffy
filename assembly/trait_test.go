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

func writeTraits(db *sql.DB, pairs []pair) (err error) {
	dbq := ephemera.NewDBQueue(db)
	r := ephemera.NewRecorder("test", dbq)
	w := NewWriter(dbq)

	for _, p := range pairs {
		aspect := r.Named(ephemera.NAMED_ASPECT, p.key, "key")
		trait := r.Named(ephemera.NAMED_TRAIT, p.value, "value")
		r.NewAspect(aspect)
		r.NewTrait(trait, aspect, 0)
	}
	if e := DetermineTraits(w, db); e != nil {
		err = e
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

func TestTraits(t *testing.T) {
	const source = "file:test.db?cache=shared&mode=memory"
	if db, e := sql.Open("sqlite3", source); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		if e := writeTraits(db,
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
			t.Fatal(e)
		}
	}
}
