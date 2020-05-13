package assembly

import (
	"database/sql"
	"strconv"
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
)

func TestPatternCheck(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()

		// okay.
		addEphPattern(t.rec,
			"pat", "num", "number_eval", "0",
			"pat", "", "number_eval", "0",
			"pat", "", "number_eval", "1",
			"pat", "num", "number_eval", "1",
		)
		if e := checkPatternSetup(t.db); e != nil {
			t.Fatal(e)
		} else {
			t.Log("ok", "normal pattern declartion and usage")
		}
		tables.Must(t.db, "delete from eph_pattern")

		// okay.
		addEphPattern(t.rec,
			"num", "", "number_eval", "1",
			"txt", "", "text_eval", "1",
			"exe", "", "execute", "1",
		)
		if e := checkPatternSetup(t.db); e != nil {
			t.Fatal(e)
		} else {
			t.Log("ok", "three completely different pattern decls")
		}
		tables.Must(t.db, "delete from eph_pattern")

		// never declared return
		addEphPattern(t.rec,
			"pat", "num", "number_eval", "0",
			"pat", "num", "number_eval", "1",
		)
		if e := checkPatternSetup(t.db); e != nil {
			t.Log("ok", e)
		} else {
			t.Fatal("expected never declared return")
		}
		tables.Must(t.db, "delete from eph_pattern")

		// referenced an undeclared arg
		addEphPattern(t.rec,
			"pat", "", "number_eval", "1",
			"pat", "num", "number_eval", "0",
		)
		if e := checkPatternSetup(t.db); e != nil {
			t.Log("ok", e)
		} else {
			t.Fatal("expected undeclared arg")
		}
		tables.Must(t.db, "delete from eph_pattern")

		// referenced an undeclared pattern
		addEphPattern(t.rec,
			"pat", "", "number_eval", "0",
		)
		if e := checkPatternSetup(t.db); e != nil {
			t.Log("ok", e)
		} else {
			t.Fatal("expected undeclared pat")
		}
		tables.Must(t.db, "delete from eph_pattern")

		// arg mismatch
		addEphPattern(t.rec,
			"pat", "", "number_eval", "1",
			"pat", "num", "number_eval", "1",
			"pat", "num", "text_eval", "1",
		)
		if e := checkPatternSetup(t.db); e != nil {
			t.Log("ok", e)
		} else {
			t.Fatal("expected type mismatch")
		}
		tables.Must(t.db, "delete from eph_pattern")

		// return mismatch
		addEphPattern(t.rec,
			"pat", "", "number_eval", "1",
			"pat", "", "text_eval", "1",
		)
		if e := checkPatternSetup(t.db); e != nil {
			t.Log("ok", e)
		} else {
			t.Fatal("expected type mismatch")
		}
		tables.Must(t.db, "delete from eph_pattern")

		// variable and pattern names in the same pattern shouldnt match
		addEphPattern(t.rec,
			"pat", "pat", "number_eval", "1",
		)
		if e := checkPatternSetup(t.db); e != nil {
			t.Log("ok", e)
		} else {
			t.Fatal("expected name conflict")
		}
		tables.Must(t.db, "delete from eph_pattern")

		// variable and pattern names in the same pattern shouldnt match
		addEphPattern(t.rec,
			"pat", "", "text_eval", "1",
			"pat", "bat", "number_eval", "1",
			//
			"bat", "", "text_eval", "1",
			"bat", "pat", "number_eval", "1",
		)
		if e := checkPatternSetup(t.db); e != nil {
			t.Fatal(e)
		} else {
			t.Log("ok", "variable and pattern names can be reused")
		}
		tables.Must(t.db, "delete from eph_pattern")
	}
}

func checkPatternSetup(db *sql.DB) (err error) {
	var now, last patternInfo
	var declaredReturn string

	// find where variable names and pattern names conflict
	if e := tables.QueryAll(db,
		`select distinct pn.name 
		from eph_pattern ep
		left join eph_named pn
			on (ep.idNamedPattern = pn.rowid)
		left join eph_named kn
			on (ep.idNamedParam = kn.rowid)
		where ep.idNamedPattern != ep.idNamedParam
		and pn.name = kn.name`,
		func() error {
			e := now.compare(&last, &declaredReturn)
			last = now
			return e
		},
		&now.pat, &now.arg, &now.typ, &now.decl); e != nil {
		err = e
	} else {
		// search for other conflicts
		if e := tables.QueryAll(db,
			`select distinct * from asm_pattern
			order by pattern, param, type, decl desc`,
			func() error {
				e := now.compare(&last, &declaredReturn)
				last = now
				return e
			},
			&now.pat, &now.arg, &now.typ, &now.decl); e != nil {
			err = e
		} else if e := last.flush(&declaredReturn); e != nil {
			err = e
		}
	}
	return
}

type patternInfo struct {
	pat, arg, typ string
	decl          bool
}

func (now *patternInfo) flush(pret *string) (err error) {
	if len(now.pat) > 0 && len(*pret) == 0 {
		err = errutil.Fmt("Pattern %q never declared a return type", now.pat)
	}
	(*pret) = ""
	return
}

func (now *patternInfo) compare(was *patternInfo, pret *string) (err error) {
	// new pattern detected....
	if now.pat != was.pat {
		if e := was.flush(pret); e != nil {
			err = e
		}
	}
	if err == nil {
		if change := (now.pat != was.pat || now.arg != was.arg); change && !now.decl {
			// decl(s) come first, so if there's a change... it should only happen with a decl.
			err = errutil.Fmt("Pattern %q's %q missing declaration", now.pat, now.arg)
		} else if !change && (now.typ != was.typ) {
			// regardless -- types should be consistent.
			err = errutil.New("Pattern %q's %q type conflict, was %q now %q", now.pat, now.arg, was.typ, now.typ)
		} else if now.decl && now.pat == now.arg {
			// assuming everything's ok, a decl where pat and arg match means the type of the pattern itself.
			*pret = now.typ
		}
	}
	return
}

// adds rows of 4 values to the database of test ephemera
func addEphPattern(rec *ephemera.Recorder, els ...string) {
	for i := 0; i < len(els); i += 4 {
		pat := rec.NewName(els[i+0], tables.NAMED_PATTERN, strconv.Itoa(i))
		arg := pat
		if n := els[i+1]; len(n) > 0 {
			arg = rec.NewName(els[i+1], tables.NAMED_VARIABLE, strconv.Itoa(i))
		}
		typ := rec.NewName(els[i+2], tables.NAMED_TYPE, strconv.Itoa(i))

		if dec, _ := strconv.ParseBool(els[i+3]); dec {
			rec.NewPatternDecl(pat, arg, typ)
		} else {
			rec.NewPatternRef(pat, arg, typ)
		}
	}
}
