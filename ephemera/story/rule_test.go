package story_test

import (
	"strings"
	"testing"

	"github.com/ionous/iffy/ephemera/story"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/print"
	"github.com/ionous/iffy/rt/safe"
	"github.com/ionous/iffy/rt/writer"
	"github.com/ionous/iffy/tables"
	"github.com/ionous/iffy/test/testdb"
	"github.com/ionous/iffy/test/testutil"
	"github.com/kr/pretty"
)

func TestPatternActivity(t *testing.T) {
	_, decoder, db := newImporter(t, testdb.Memory)
	defer db.Close()
	if prog, e := decoder.ReadSpec(_pattern_activity); e != nil {
		t.Fatal(e)
	} else if exe, ok := prog.(rt.Execute); !ok {
		t.Fatalf("cant cast %T to execute", exe)
	} else {
		var run testRuntime
		out := print.NewLines()
		run.SetWriter(out)
		// should this call/test buildRule
		if e := safe.Run(&run, exe); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out.Lines(), []string{"hello", "hello"}); len(diff) > 0 {
			t.Fatal(diff)
		}
	}
}

func TestPatternRule(t *testing.T) {
	k, decoder, db := newImporter(t, testdb.Memory)
	defer db.Close()
	if i, e := decoder.ReadSpec(_pattern_actions); e != nil {
		t.Fatal(e)
	} else if act, ok := i.(*story.PatternActions); !ok {
		t.Fatalf("cant cast %T to pattern actions", i)
	} else if e := act.ImportPhrase(k); e != nil {
		t.Fatal(e)
	} else {
		var buf strings.Builder
		// -- eph?rule
		// execute_rule
		tables.WriteCsv(db, &buf, "select progType from eph_prog", 1)
		// example, pattern_name
		tables.WriteCsv(db, &buf, "select name, category from eph_named where category != 'scene'", 2)
		// 1 pattern handler reference
		tables.WriteCsv(db, &buf, "select count() from eph_pattern", 1)
		// 1, 1 - the first name, the first program are used to make the rule
		tables.WriteCsv(db, &buf, "select idNamedPattern, idProg from eph_rule", 2)
		if have, want := buf.String(), lines(
			"execute_rule",
			"example,pattern",
			"0", // eph_pattern -- rules are recorded via eph_prog,
			"2,1",
		); have != want {
			t.Fatal(have)
		}
	}
}

var _pattern_actions = map[string]interface{}{
	"type": "pattern_actions",
	"value": map[string]interface{}{
		"$NAME": map[string]interface{}{
			"type":  "pattern_name",
			"value": "example",
		},
		"$PATTERN_RULES": map[string]interface{}{
			"type": "pattern_rules",
			"value": map[string]interface{}{
				"$PATTERN_RULE": []interface{}{
					map[string]interface{}{
						"type": "pattern_rule",
						"value": map[string]interface{}{
							"$GUARD": map[string]interface{}{
								"type": "bool_eval",
								"value": map[string]interface{}{
									"type":  "always",
									"value": map[string]interface{}{},
								}},

							"$HOOK": map[string]interface{}{
								"type": "program_hook",
								"value": map[string]interface{}{
									"$ACTIVITY": _pattern_activity,
								}}}}}}}},
}

var _pattern_activity = map[string]interface{}{
	"type": "activity",
	"value": map[string]interface{}{
		"$EXE": []interface{}{
			_say_exec,
			_say_exec,
		},
	},
}

var _say_exec = map[string]interface{}{
	"type": "execute",
	"value": map[string]interface{}{
		"type": "say_text",
		"value": map[string]interface{}{
			"$TEXT": _text_eval,
		},
	},
}

var _text_eval = map[string]interface{}{
	"type": "text_eval",
	"value": map[string]interface{}{
		"type": "text_value",
		"value": map[string]interface{}{
			"$TEXT": map[string]interface{}{
				"type":  "text",
				"value": "hello",
			},
		},
	},
}

type baseRuntime struct {
	testutil.PanicRuntime
}
type testRuntime struct {
	baseRuntime
	writer.Sink
}
