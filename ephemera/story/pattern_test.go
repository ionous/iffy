package story

import (
	"strings"
	"testing"

	"github.com/ionous/iffy/tables"
	"github.com/ionous/iffy/test/testdb"
)

func TestPatternVars(t *testing.T) {
	patternVariables := &PatternVariablesDecl{
		PatternName: PatternName{
			Str: "corral",
		},
		VariableDecl: []VariableDecl{{
			Type: VariableType{
				Opt: &ObjectType{
					An: An{
						Str: "$AN",
					},
					Kind: SingularKind{
						Str: "animal",
					},
				},
			},
			Name: VariableName{
				Str: "pet",
			},
		},
		},
	}
	k, db := newImporter(t, testdb.Memory)
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
	patternDecl := &PatternDecl{
		Name: PatternName{
			Str: "corral",
		},
		Type: PatternType{
			Opt: &PatternedActivity{
				Str: "$ACTIVITY",
			},
		},
	}

	k, db := newImporter(t, testdb.Memory)
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
