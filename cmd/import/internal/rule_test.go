package internal

import (
	"strings"
	"testing"

	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/print"
	"github.com/ionous/iffy/tables"
	"github.com/kr/pretty"
)

// import an object type description
func TestObjectFunc(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if rule, e := imp_object_func(k, _object_func); e != nil {
		t.Fatal(e)
	} else if text, e := rt.GetText(nil, rule.distill().(rt.TextEval)); e != nil {
		t.Fatal(e)
	} else if text != "hello" {
		t.Fatal(text)
	}
}

func TestPatternActivity(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if rule, e := imp_pattern_activity(k, _pattern_activity); e != nil {
		t.Fatal(e)
	} else {
		var run testRuntime
		if e := rt.RunOne(&run, rule.distill().(rt.Execute)); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(run.out.Lines(), []string{"hello", "hello"}); len(diff) > 0 {
			t.Fatal(diff)
		}
	}
}

func TestPatternHandler(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if e := imp_pattern_handler(k, _pattern_handler); e != nil {
		t.Fatal(e)
	} else {
		var buf strings.Builder
		// exec_rule
		tables.WriteCsv(db, &buf, "select type from eph_prog", 1)
		// example, pattern_name
		tables.WriteCsv(db, &buf, "select name, category from eph_named", 2)
		// 0 -- we dont have the pattern definition, just the rules
		tables.WriteCsv(db, &buf, "select count() from eph_pattern", 1)
		// 1, 1 - the first name, the first program are used to make the rule
		tables.WriteCsv(db, &buf, "select idNamedPattern, idProg from eph_rule", 2)
		if diff := pretty.Diff(buf.String(), lines(
			"exec_rule",
			"example,pattern_name",
			"0",
			"1,1",
		)); len(diff) > 0 {
			t.Fatal(diff)
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
	"type": "pattern_activity",
	"value": map[string]interface{}{
		"$GO": []interface{}{
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
				"type":  "lines",
				"value": "hello",
			},
		},
	},
}

type testRuntime struct {
	rt.Panic
	out print.Lines
}

func (t *testRuntime) Write(p []byte) (ret int, err error) {
	return t.out.Write(p)
}
