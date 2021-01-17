package story_test

import (
	"strings"
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ephemera/story"
	"github.com/ionous/iffy/tables"
	"github.com/ionous/iffy/test/testdb"
)

func TestPatternVars(t *testing.T) {
	patternVariables := &story.PatternVariablesDecl{
		PatternName: story.PatternName{
			Str: "corral",
		},
		VariableDecl: []story.VariableDecl{{
			Type: story.VariableType{
				Opt: &story.ObjectType{
					An: story.Ana{
						Str: "$AN",
					},
					Kind: story.SingularKind{
						Str: "animal",
					},
				},
			},
			Name: story.VariableName{core.Variable{
				Str: "pet",
			}},
		}},
	}
	k, _, db := newImporter(t, testdb.Memory)
	defer db.Close()

	if e := patternVariables.ImportPhrase(k); e != nil {
		t.Log(e)
	} else {
		var buf strings.Builder
		tables.WriteCsv(db, &buf, "select name, category from eph_named where category != 'scene'", 2)
		tables.WriteCsv(db, &buf, "select idNamedPattern,idNamedParam,idNamedType,idProg from eph_pattern", 4)
		if have, want := buf.String(), lines(
			"corral,pattern",       // 1
			"pet,parameter",        // 2
			"animal,singular_kind", // 3
			"2,3,4,0",              // NewPatternDecl
		); have != want {
			t.Fatal("mismatch, have:", have)
		} else {
			t.Log("ok")
		}
	}
}

func TestPatternDecl(t *testing.T) {
	patternDecl := &story.PatternDecl{
		Name: story.PatternName{
			Str: "corral",
		},
		Type: story.PatternType{
			Opt: &story.PatternedActivity{
				Str: "$ACTIVITY",
			},
		},
	}

	k, _, db := newImporter(t, testdb.Memory)
	defer db.Close()
	if e := patternDecl.ImportPhrase(k); e != nil {
		t.Fatal(e)
	} else {
		var buf strings.Builder
		tables.WriteCsv(db, &buf, "select name, category from eph_named where category != 'scene'", 2)
		tables.WriteCsv(db, &buf, "select idNamedPattern,idNamedParam,idNamedType,idProg from eph_pattern", 4)
		if have, want := buf.String(), lines(
			"corral,pattern", // 1
			"execute,type",   // 2
			"2,2,3,0",        // NewPatternDecl
		); have != want {
			t.Fatal("mismatch", have)
		} else {
			t.Log("ok")
		}
	}
}
