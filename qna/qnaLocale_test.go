package qna

import (
	"strings"
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/rel"
	"github.com/ionous/iffy/rt/safe"
	"github.com/ionous/iffy/tables"
	"github.com/ionous/iffy/test/testdb"
)

//
func TestLocale(t *testing.T) {
	if db, e := newPairTest(t, testdb.Memory); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		run_domain := testdb.TableCols("run_domain", "domain", "active")
		mdl_kind := testdb.TableCols("mdl_kind", "kind", "path")
		mdl_noun := testdb.TableCols("mdl_noun", "noun", "kind")
		mdl_rel := testdb.TableCols("mdl_rel", "relation", "cardinality", "kind", "otherKind")
		mdl_name := testdb.TableCols("mdl_name", "noun", "name")
		run_pair := testdb.TableCols("run_pair", "relation", "noun", "otherNoun")
		const p, q string = "#test::p", "#test::q"

		if e := testdb.Ins(db, run_domain,
			"test", 1,
		); e != nil {
			t.Fatal(e)
		} else if e := testdb.Ins(db, mdl_kind,
			"Ts", "",
		); e != nil {
			t.Fatal(e)
		} else if e := testdb.Ins(db, mdl_noun,
			p, "Ts",
			q, "Ts",
		); e != nil {
			t.Fatal(e)
		} else if e := testdb.Ins(db, mdl_name,
			p, "p",
			q, "q",
		); e != nil {
			t.Fatal(e)
		} else if e := testdb.Ins(db, mdl_rel,
			"locale", tables.MANY_TO_ONE, "Ts", "Ts",
		); e != nil {
			t.Fatal(e)
		} else {
			run := NewRuntime(db)
			if e := safe.Run(run, &rel.Reparent{&core.Text{"q"}, &core.Text{"p"}}); e != nil {
				t.Fatal(e)
			} else if v, e := safe.GetObject(run, &rel.Locale{&core.ObjectName{&core.Text{"q"}}}); e != nil {
				t.Fatal(e)
			} else if v.String() != p {
				t.Fatal("unexpected parent", v.String())
			} else {
				var buf strings.Builder
				if e := testdb.WriteCsv(db, &buf, run_pair, "where active=1"); e != nil {
					t.Fatal(e)
				} else if res, want := buf.String(), "locale,#test::q,#test::p"; res != lines(want) {
					t.Log("locale not written to db", res)
				}
			}
		}
	}
}
