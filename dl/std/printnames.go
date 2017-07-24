package std

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/std/group"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
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
		var buffer printer.Span
		run.PushWriter(printer.AndSeparator(&buffer))
		//
		if e := printWithArticles(run, ungrouped); e != nil {
			err = e
		} else if e := printGroups(run, groups); e != nil {
			err = e
		}
		run.PopWriter()
		if err == nil {
			_, e := run.Write(buffer.Bytes())
			err = e
		}
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
		_, err = run.ExecuteMatching(printName)
	}
	return
}
