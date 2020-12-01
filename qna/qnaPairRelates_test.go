package qna

import (
	"testing"

	"github.com/ionous/iffy/tables"
	"github.com/ionous/iffy/test/testdb"
)

func TestPairRelates(t *testing.T) {
	if db, e := newPairTest(t, testdb.Memory); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		run_domain := testdb.TableCols("run_domain", "domain", "active")
		mdl_kind := testdb.TableCols("mdl_kind", "kind", "path")
		mdl_noun := testdb.TableCols("mdl_noun", "noun", "kind")
		mdl_rel := testdb.TableCols("mdl_rel", "relation", "cardinality", "kind", "otherKind")
		run_pair := testdb.TableCols("run_pair", "relation", "noun", "otherNoun")
		const rel, p, q, r string = "R", "#test::p", "#test::q", "#test::r"

		if e := testdb.Ins(db, run_domain,
			"test", 1,
		); e != nil {
			t.Fatal(e)
		} else if e := testdb.Ins(db, mdl_kind,
			"Ts", "",
			"Ps", "Ts",
			"Qs", "Ts",
		); e != nil {
			t.Fatal(e)
		} else if e := testdb.Ins(db, mdl_noun,
			p, "Ps",
			q, "Qs",
			r, "Qs",
		); e != nil {
			t.Fatal(e)
		} else if e := testdb.Ins(db, mdl_rel,
			rel, tables.ONE_TO_MANY, "Qs", "Ps",
		); e != nil {
			t.Fatal(e)
		} else if e := testdb.Ins(db, run_pair,
			rel, q, p,
		); e != nil {
			t.Fatal(e)
		} else {
			run := NewRuntime(db)
			if vs, e := run.RelativesOf(q, rel); e != nil {
				t.Fatal(e)
			} else if len(vs) != 1 || vs[0] != p {
				t.Fatal("initial relatives mismatched", vs)
			} else if e := run.RelateTo(r, p, rel); e != nil {
				t.Fatal(e)
			} else if vs, e := run.RelativesOf(q, rel); e != nil {
				t.Fatal(e)
			} else if len(vs) != 0 {
				t.Fatal("original relatives not cleared", vs)
			} else if vs, e := run.RelativesOf(r, rel); e != nil {
				t.Fatal(e)
			} else if len(vs) != 1 || vs[0] != p {
				t.Fatal("related relatives mismatched", vs)
			}
		}
	}
}
