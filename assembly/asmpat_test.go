package assembly

import (
	"strings"
	"testing"

	"github.com/ionous/iffy/tables"
)

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
