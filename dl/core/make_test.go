package core

import (
	"testing"

	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/test"
	"github.com/kr/pretty"
)

func TestMake(t *testing.T) {
	type panicTime struct {
		rt.Panic
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
	} else if d, e := obj.GetRecord(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(g.RecordToValue(d), map[string]interface{}{
		"Name":         "",
		"Label":        "",
		"Innumerable":  "NotInnumerable",
		"GroupOptions": "WithArticles",
	}); len(diff) != 0 {
		t.Fatal(diff)
	}
}
