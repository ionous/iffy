package story

import (
	"strings"
	"testing"

	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/print"
	"github.com/ionous/iffy/rt/writer"
	"github.com/ionous/iffy/tables"
	"github.com/kr/pretty"
)

// import an object type description
func TestObjectFunc(t *testing.T) {
	k, db := newTestDecoder(t)
	defer db.Close()
	if rule, e := imp_object_func(k, _object_func); e != nil {
		t.Fatal(e)
	} else if text, e := rt.GetText(nil, rule.buildRule().(rt.TextEval)); e != nil {
		t.Fatal(e)
	} else if text != "hello" {
		t.Fatal(text)
	}
}

func TestPatternActivity(t *testing.T) {
	k, db := newTestDecoder(t)
	defer db.Close()
	if rule, e := imp_activity(k, _pattern_activity); e != nil {
		t.Fatal(e)
	} else {
		var run testRuntime
		out := print.NewLines()
		run.SetWriter(out)
		// should this call/test buildRule
		if e := rt.RunOne(&run, rule); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out.Lines(), []string{"hello", "hello"}); len(diff) > 0 {
			t.Fatal(diff)
		}
	}
}

func TestPatternRule(t *testing.T) {
	k, db := newTestDecoder(t)
	defer db.Close()
	if e := imp_pattern_handler(k, _pattern_handler); e != nil {
		t.Fatal(e)
	} else {
		var buf strings.Builder
		// execute_rule
		tables.WriteCsv(db, &buf, "select type from eph_prog", 1)
		// example, pattern_name
		tables.WriteCsv(db, &buf, "select name, category from eph_named where category != 'scene'", 2)
		// 1 pattern handler reference
		tables.WriteCsv(db, &buf, "select count() from eph_pattern", 1)
		// 1, 1 - the first name, the first program are used to make the rule
		tables.WriteCsv(db, &buf, "select idNamedPattern, idProg from eph_rule", 2)
		if have, want := buf.String(), lines(
			"execute_rule",
			"example,pattern_name",
			"0", // eph_pattern -- rules are recorded via eph_prog,
			"2,1",
		); have != want {
			t.Fatal(have)
		}
	}
}

var _pattern_handler = map[string]interface{}{
	"type": "pattern_handler",
	"value": map[string]interface{}{
		"$NAME": map[string]interface{}{
			"type":  "pattern_name",
			"value": "example",
		},
		"$HOOK": map[string]interface{}{
			"type": "pattern_hook",
			"value": map[string]interface{}{
				"$ACTIVITY": _pattern_activity,
			}}},
}

var _object_func = map[string]interface{}{
	"type":  "object_func",
	"value": _text_eval,
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
	rt.Panic
}
type testRuntime struct {
	baseRuntime
	writer.Sink
}
