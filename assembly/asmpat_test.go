package assembly

import (
	"strings"
	"testing"

	"github.com/ionous/iffy/tables"
)

//
func TestPatternAsm(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		t.Run("normal", func(t *testing.T) {
			cleanPatterns(asm.db)
			// name, param, type, decl
			addEphPattern(asm.rec,
				"put", "z", "text_eval", "true",
				"pat", "", "number_eval", "true",
				"pat", "a", "number_eval", "true",
				"pat", "b", "number_eval", "true",
				"pat", "b", "number_eval", "false",
				"pat", "c", "number_eval", "true",
				"pat", "w", "number_eval", "false",
				"put", "", "text_eval", "true",
			)
			if e := copyPatterns(asm.db); e != nil {
				t.Fatal(e)
			} else {

				var buf strings.Builder
				tables.WriteCsv(asm.db, &buf, "select * from mdl_pat", 4)
				if have, want := buf.String(), lines(
					"pat,pat,number_eval,0",
					"pat,a,number_eval,1",
					"pat,b,number_eval,2",
					"pat,c,number_eval,3",
					//
					"put,put,text_eval,0",
					"put,z,text_eval,1",
				); have != want {
					t.Fatal(have)
				} else {
					t.Log("ok")
				}
			}
		})
	}
}

// assemble patterns and rules into the model rules and programs.
// doesnt check for consistency -- that's up to pattern and rule checking currently.
func TestRuleAsm(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		t.Run("normal", func(t *testing.T) {
			cleanPatterns(asm.db)
			// name, param, type, decl
			addEphPattern(asm.rec,
				"pat", "", "number_eval", "true",
				"put", "", "text_eval", "true")
			addEphRule(asm.rec,
				"pat", "number_rule", "some number",
				"put", "text_rule", "other text",
				"pat", "number_rule", "other number",
			)
			if e := copyRules(asm.db); e != nil {
				t.Fatal(e)
			} else {
				var buf strings.Builder
				tables.WriteCsv(asm.db, &buf, "select count() from mdl_rule", 1)
				if have, want := buf.String(), lines(
					"3",
				); have != want {
					t.Fatal(have)
				} else {
					t.Log("ok")
				}
			}
		})
	}
}

func lines(s ...string) string {
	return strings.Join(s, "\n") + "\n"
}
