package std

import (
	"github.com/ionous/iffy/dl/std/group"
	"github.com/ionous/iffy/ref/obj"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"github.com/ionous/iffy/rt/stream"
)

// PrintNondescriptObjects prints a bunch of objects, in groups if possible,
// using the GroupTogether, PrintGroup, and PrintName patterns.
// This is similar to the Inform activity "Listing contents of something" and its I6 "standard contents listing rule".
// http://inform7.com/learn/man/WI_18_13.html
type PrintNondescriptObjects struct {
	Objects rt.ObjListEval
}

func (p *PrintNondescriptObjects) Execute(run rt.Runtime) (err error) {
	if groups, ungrouped, e := group.MakeGroups(run, p.Objects); e != nil {
		err = e
	} else {
		sep := printer.AndSeparator(run)
		run, ungrouped := rt.Writer(run, sep), stream.NewObjectStream(ungrouped)
		if _, e := printWithArticles(run, ungrouped); e != nil {
			err = e
		} else if e := groups.PrintGroups(run); e != nil {
			err = e
		}
		sep.Flush()
	}
	return
}

func printName(run rt.Runtime, x rt.Object) error {
	printName := obj.Emplace(&PrintName{x.Id()})
	return run.ExecuteMatching(run, printName)
}

func printPluralName(run rt.Runtime, x rt.Object) error {
	printName := obj.Emplace(&PrintPluralName{x.Id()})
	return run.ExecuteMatching(run, printName)
}

func printSeveral(run rt.Runtime, x rt.Object, cnt int) error {
	printName := obj.Emplace(&PrintSeveral{x.Id(), float64(cnt)})
	return run.ExecuteMatching(run, printName)
}
