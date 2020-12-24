package test

import (
	"testing"

	"github.com/ionous/iffy/dl/core"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/test/testutil"
	"github.com/kr/pretty"
)

func TestMake(t *testing.T) {
	type panicTime struct {
		testutil.PanicRuntime
	}
	var testTime struct {
		panicTime
		testutil.Kinds
	}
	testTime.Kinds.AddKinds((*GroupSettings)(nil))
	op := &core.Make{Name: "GroupSettings",
		Arguments: &core.Arguments{[]*core.Argument{
			{"ObjectsWithArticles", &core.FromBool{&core.Bool{true}}},
		}}}
	if obj, e := op.GetRecord(&testTime); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(g.RecordToValue(obj.Record()), map[string]interface{}{
		"Name":         "",
		"Label":        "",
		"Innumerable":  "NotInnumerable",
		"GroupOptions": "ObjectsWithArticles",
	}); len(diff) != 0 {
		t.Fatal(diff)
	}
}
