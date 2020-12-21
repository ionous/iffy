package assembly

import (
	"encoding/gob"
	"strconv"
	"strings"
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/ephemera/debug"
	"github.com/ionous/iffy/tables"
	"github.com/ionous/iffy/test/testdb"
	"github.com/kr/pretty"
)

//
func TestPatternAsm(t *testing.T) {
	if asm, e := newAssemblyTest(t, testdb.Memory); e != nil {
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
	gob.Register((*core.Text)(nil))
	gob.Register((*debug.MatchNumber)(nil))
	//
	if asm, e := newAssemblyTest(t, testdb.Memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		t.Run("normal", func(t *testing.T) {
			cleanPatterns(asm.db)
			// say_me return text, and has one number parameter
			addEphPattern(asm.rec,
				"say_me", "", "text_eval", "true",
				"say_me", "num", "number_eval", "true")

			for i, rule := range debug.SayPattern.Rules {
				prog, e := asm.rec.NewGob("text_rule", &rule)
				if e != nil {
					t.Fatal(e)
				}
				asm.rec.NewPatternRule(
					asm.rec.NewName("say_me",
						tables.NAMED_PATTERN,
						strconv.Itoa(i+1)),
					prog)
			}
			//
			if e := buildPatterns(asm.assembler); e != nil {
				t.Fatal(e)
			} else {
				var visited bool
				var progName, typeName string
				var pat pattern.TextPattern
				if e := tables.QueryAll(asm.db, "select * from mdl_prog",
					func() (err error) {
						if visited {
							err = errutil.New("multiple programs detected")
						}
						visited = true
						return
					}, &progName, &typeName, tables.NewGobScanner(&pat)); e != nil {
					t.Fatal(e)
				} else if progName != "say_me" || typeName != "TextPattern" {
					t.Fatal("mismatched columns")
				} else if diff := pretty.Diff(debug.SayPattern, pat); len(diff) > 0 {
					pretty.Println("want:", debug.SayPattern)
					pretty.Println("have:", pat)
					t.Fatal("error")
				}
			}
		})
	}
}

func lines(s ...string) string {
	return strings.Join(s, "\n") + "\n"
}
