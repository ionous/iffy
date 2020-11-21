package core

import (
	"testing"

	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/test"
	"github.com/kr/pretty"
)

func TestMake(t *testing.T) {
	type panicTime struct {
		test.PanicRuntime
	}
	var testTime struct {
		panicTime
		test.Kinds
	}
	testTime.Kinds.AddKinds((*test.GroupSettings)(nil))
	op := Make{Name: "GroupSettings",
		Arguments: &Arguments{[]*Argument{
			{"WithArticles", &FromBool{&Bool{true}}},
		}}}
	if obj, e := op.GetObject(&testTime); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(g.RecordToValue(obj.Record()), map[string]interface{}{
		"Name":         "",
		"Label":        "",
		"Innumerable":  "NotInnumerable",
		"GroupOptions": "WithArticles",
	}); len(diff) != 0 {
		t.Fatal(diff)
	}
}
