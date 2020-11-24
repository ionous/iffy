package test

import (
	"github.com/ionous/iffy/object"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/scope"
	"github.com/ionous/iffy/rt/writer"
	"github.com/ionous/iffy/test/testutil"
)

type panicTime struct {
	testutil.PanicRuntime
}

type testTime struct {
	panicTime
	objs map[string]*g.Record
	scope.ScopeStack
	testutil.PatternMap
	*testutil.Kinds
}

func (lt *testTime) Writer() writer.Output {
	return writer.NewStdout()
}

func (lt *testTime) GetField(target, field string) (ret g.Value, err error) {
	if obj, ok := lt.objs[field]; target == object.Value && ok {
		ret = g.RecordOf(obj)
	} else {
		ret, err = lt.ScopeStack.GetField(target, field)
	}
	return
}
