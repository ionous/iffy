package qna

import (
	"database/sql"
	"log"
	"strings"
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/tables"
	"github.com/ionous/iffy/test/testdb"
)

func TestPairChanges(t *testing.T) {
	if db, e := newPairTest(t, testdb.Memory); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		if fields, e := NewFields(db); e != nil {
			t.Fatal(e)
		} else {
			if e := testdb.Ins(db, testdb.TableCols("mdl_rel", "relation", "cardinality" /*kind, otherKind*/),
				"11", tables.ONE_TO_ONE,
				"1N", tables.ONE_TO_MANY,
				"N1", tables.MANY_TO_ONE,
				"NN", tables.MANY_TO_MANY,
			); e != nil {
				t.Fatal(e)
			}
			// existing relations
			run_pair := testdb.TableCols("run_pair", "relation", "noun", "otherNoun")
			if e := testdb.Ins(db, run_pair,
				// a noun has at most one other noun which refers to it and only it.
				"11", "a", "b",
				"11", "c", "d",
				"11", "n", "m",
				// a noun has at most one other noun which refers to it.
				"1N", "a", "b",
				"1N", "a", "c",
				// a number of nouns each have a single reference to some other noun.
				"N1", "a", "b",
				"N1", "c", "b",
				// various nouns reference various other nouns.
				"NN", "a", "b",
				"NN", "b", "a",
			); e != nil {
				t.Fatal(e)
			}
			// prepare to update those relations
			if e := testdb.Ins(db, testdb.TableCols("mdl_pair", "relation", "noun", "otherNoun"),
				// should replace both the a-b, and the c-d relations.
				"11", "a", "d",
				// should replace the a-c pair; a-b should be replace itself.
				"1N", "a", "b",
				"1N", "d", "c",
				// should add the d-b relation, should replace the a-b relation
				"N1", "d", "b",
				"N1", "a", "g",
				// should add the a-f and f-a relations, everything else should be unchanged
				"NN", "a", "f",
				"NN", "f", "a",
				"NN", "a", "b",
			); e != nil {
				t.Fatal(e)
			}
			// run the test:
			if cnt, e := fields.UpdatePairs("entireGame"); e != nil {
				t.Fatal(e)
			} else {
				log.Println("activate domain affected", cnt, "rows")
				var buf strings.Builder
				if e := testdb.WriteCsv(db, &buf, run_pair); e != nil {
					t.Fatal(e)
				} else if have, want := buf.String(), lines(
					"11,a,d",
					"11,n,m",
					"1N,a,b",
					"1N,d,c",
					"N1,a,g",
					"N1,c,b",
					"N1,d,b",
					"NN,a,b",
					"NN,a,f",
					"NN,b,a",
					"NN,f,a",
				); have != want {
					t.Fatal(have)
				}
			}
		}
	}
}
func lines(s ...string) string {
	return strings.Join(s, "\n") + "\n"
}

func newPairTest(t *testing.T, path string) (ret *sql.DB, err error) {
	var source string
	if len(path) > 0 {
		source = path
	} else if p, e := testdb.PathFromName(t.Name()); e != nil {
		t.Fatal(e)
	} else {
		source = p
	}
	if db, e := sql.Open(tables.DefaultDriver, source); e != nil {
		err = errutil.New(e, "for", source)
	} else if e := tables.CreateModel(db); e != nil {
		err = errutil.New(e, "for", source)
	} else if e := tables.CreateRun(db); e != nil {
		err = errutil.New(e, "for", source)
	} else if e := tables.CreateRunViews(db); e != nil {
		err = errutil.New(e, "for", source)
	} else {
		ret = db
	}
	return
}
