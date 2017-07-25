package std

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/std/group"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"github.com/ionous/iffy/rt/stream"
	r "reflect"
)

// PrintNondescriptObjects commands the runtime to print a bunch of objects, in groups if possible.
type PrintNondescriptObjects struct {
	Objects rt.ObjListEval
}

func (p *PrintNondescriptObjects) Execute(run rt.Runtime) (err error) {
	if groups, ungrouped, e := group.MakeGroups(run, p.Objects); e != nil {
		err = e
	} else {
		sep := printer.AndSeparator(run)
		run, ungrouped := rt.Writer(run, sep), stream.NewObjectStream(ungrouped)
		if e := printWithArticles(run, ungrouped); e != nil {
			err = e
		} else if e := groups.Print(run); e != nil {
			err = e
		}
		sep.Flush()
	}
	return
}

// FIX: this is patently ridiculous.
// issue: i cant set an object reference from an object
// why? in part b/c theres no "base class"
// it would be **alot** simpler if the * was an ident.Id
// we'd still have "emplace" -- you could maybe someday make it static -- thatd be tons better.
func printName(run rt.Runtime, obj rt.Object) (err error) {
	var kind *Kind
	if src, ok := obj.(*ref.RefObject); !ok {
		err = errutil.Fmt("unknown object %T", obj)
	} else if e := ref.Upcast(src.Value().Addr(), func(ptr r.Value) (okay bool) {
		kind, okay = ptr.Interface().(*Kind)
		return
	}); e != nil {
		err = e
	} else if printName, e := run.Emplace(&PrintName{kind}); e != nil {
		err = e
	} else {
		err = run.ExecuteMatching(run, printName)
	}
	return
}
