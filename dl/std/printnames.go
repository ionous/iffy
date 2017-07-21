package std

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	r "reflect"
)

// PrintName executes a pattern to print the target's name.
// The standard rules print the "printed name" property of the target,
// or the object name ( if the target lacks a "printed name" ),
// or the object's class name ( for unnamed objects. )
// A "printed name" can change during the course of play; object names never change.
type PrintName struct {
	Target *Kind
}

// PrintPluralName executes a pattern to print the plural of the target's name.
// The standard rules print the target's "printed plural name",
// or, if the target lacks that property, the plural of the "print name" pattern.
// It uses the runtime's pluralization table, or if needed, automated pluralization.
type PrintPluralName struct {
	Target *Kind
}

// PrintNondescriptObjects commands the runtime to print a bunch of objects, in groups if possible.
type PrintNondescriptObjects struct {
	Objects rt.ObjListEval
}

func (p *PrintNondescriptObjects) Execute(run rt.Runtime) (err error) {
	if groups, e := MakeGroups(run, p.Objects); e != nil {
		err = e
	} else {
		var buffer printer.Span
		run.PushWriter(printer.AndSeparator(&buffer))
		//
		if e := printWithArticles(run, groups.Ungrouped); e != nil {
			err = e
		} else if e := printGroups(run, groups.Grouped); e != nil {
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
