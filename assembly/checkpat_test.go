package assembly

import (
	"database/sql"
	"strconv"
	"testing"

	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
)

func TestPatternCheck(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		// okay.
		t.Run("normal", func(t *testing.T) {
			cleanPatterns(asm.db)
			addEphPattern(asm.rec,
				"pat", "num", "number_eval", "false",
				"pat", "", "number_eval", "false",
				"pat", "", "number_eval", "true",
				"pat", "num", "number_eval", "true",
			)
			if e := checkPatternSetup(asm.db); e != nil {
				t.Fatal(e)
			} else {
				t.Log("ok", "normal pattern pattern usage")
			}

		})

		// okay.
		t.Run("multi", func(t *testing.T) {
			cleanPatterns(asm.db)
			addEphPattern(asm.rec,
				"num", "", "number_eval", "true",
				"txt", "", "text_eval", "true",
				"exe", "", "execute", "true",
			)
			if e := checkPatternSetup(asm.db); e != nil {
				t.Fatal(e)
			} else {
				t.Log("ok", "three completely different pattern decls")
			}
		})

		// never declared return
		t.Run("missing return", func(t *testing.T) {
			cleanPatterns(asm.db)
			addEphPattern(asm.rec,
				"pat", "num", "number_eval", "false",
				"pat", "num", "number_eval", "true",
			)
			if e := checkPatternSetup(asm.db); e != nil {
				t.Log("ok", e)
			} else {
				t.Fatal("expected never declared return")
			}
		})

		// referenced an undeclared arg
		t.Run("undeclared arg", func(t *testing.T) {
			cleanPatterns(asm.db)
			addEphPattern(asm.rec,
				"pat", "", "number_eval", "true",
				"pat", "num", "number_eval", "false",
			)
			if e := checkPatternSetup(asm.db); e != nil {
				t.Log("ok", e)
			} else {
				t.Fatal("expected undeclared arg")
			}
		})

		// referenced an undeclared pattern
		t.Run("undeclared pattern", func(t *testing.T) {
			cleanPatterns(asm.db)
			addEphPattern(asm.rec,
				"pat", "", "number_eval", "false",
			)
			if e := checkPatternSetup(asm.db); e != nil {
				t.Log("ok", e)
			} else {
				t.Fatal("expected undeclared pat")
			}
		})
		cleanPatterns(asm.db)

		// arg mismatch
		t.Run("arg mismatch", func(t *testing.T) {
			addEphPattern(asm.rec,
				"pat", "", "number_eval", "true",
				"pat", "num", "number_eval", "true",
				"pat", "num", "text_eval", "true",
			)
			if e := checkPatternSetup(asm.db); e != nil {
				t.Log("ok", e)
			} else {
				t.Fatal("expected type mismatch")
			}
		})

		// return mismatch
		t.Run("return mismatch", func(t *testing.T) {
			cleanPatterns(asm.db)
			addEphPattern(asm.rec,
				"pat", "", "number_eval", "true",
				"pat", "", "text_eval", "true",
			)
			if e := checkPatternSetup(asm.db); e != nil {
				t.Log("ok", e)
			} else {
				t.Fatal("expected type mismatch")
			}
		})
		// variable and pattern names in the same pattern shouldnt match
		t.Run("unique variables", func(t *testing.T) {
			cleanPatterns(asm.db)
			addEphPattern(asm.rec,
				"pat", "pat", "number_eval", "true",
			)
			if e := checkPatternSetup(asm.db); e != nil {
				t.Log("ok", e)
			} else {
				t.Fatal("expected name conflict")
			}
		})
		t.Run("reused names", func(t *testing.T) {
			cleanPatterns(asm.db)
			addEphPattern(asm.rec,
				"pat", "", "text_eval", "true",
				"pat", "bat", "number_eval", "true",
				//
				"bat", "", "text_eval", "true",
				"bat", "pat", "number_eval", "true",
			)
			if e := checkPatternSetup(asm.db); e != nil {
				t.Fatal(e)
			} else {
				t.Log("ok", "variable and pattern names can be reused")
			}
		})
	}
}

// check that the pattern type matches the rule
func TestRuleCheck(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		t.Run("normal", func(t *testing.T) {
			cleanPatterns(asm.db)
			addEphPattern(asm.rec,
				"pat", "", "number_eval", "1")
			addEphRule(asm.rec,
				"pat", "number_rule", "some text")
			if e := checkRuleSetup(asm.db); e != nil {
				t.Fatal(e)
			} else {
				t.Log("ok")
			}
		})
		t.Run("multi", func(t *testing.T) {
			cleanPatterns(asm.db)
			addEphPattern(asm.rec,
				"num", "", "number_eval", "true",
				"txt", "", "text_eval", "true",
				"exe", "", "execute", "false",
			)
			addEphRule(asm.rec,
				"num", "number_rule", "some text",
				"txt", "text_rule", "some text",
			)
			if e := checkRuleSetup(asm.db); e != nil {
				t.Fatal(e)
			} else {
				t.Log("ok")
			}
		})
		//
		t.Run("missing rules", func(t *testing.T) {
			cleanPatterns(asm.db)
			addEphPattern(asm.rec,
				"pat", "", "number_eval", "1")
			if e := checkRuleSetup(asm.db); e != nil {
				t.Log("ok", e)
			} else {
				t.Error("expected missing rule")
			}
		})
		t.Run("mismatched rule type", func(t *testing.T) {
			cleanPatterns(asm.db)
			addEphPattern(asm.rec,
				"pat", "", "number_eval", "1")
			addEphRule(asm.rec,
				"pat", "text_eval", "some text")
			if e := checkRuleSetup(asm.db); e != nil {
				t.Log("ok", e)
			} else {
				t.Error("expected mismatched rule")
			}
		})
	}
}

func cleanPatterns(db *sql.DB) {
	tables.Must(db, "delete from eph_pattern")
	tables.Must(db, "delete from eph_rule")
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

// adds rows of 3 values ( name, type, text ) to the database of test ephemera
// we dont use actual "programs" here -- just strings as bytes
func addEphRule(rec *ephemera.Recorder, els ...string) {
	for i := 0; i < len(els); i += 3 {
		pat := els[i+0]
		typ := els[i+1]
		txt := els[i+2]
		rec.NewPatternRule(
			rec.NewName(pat, tables.NAMED_PATTERN, strconv.Itoa(i)),
			rec.NewProg(typ, []byte(txt)))
	}
}
