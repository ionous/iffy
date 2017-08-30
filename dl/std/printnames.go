package std

import (
	"github.com/ionous/iffy/dl/std/group"
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

func printName(run rt.Runtime, obj rt.Object) (err error) {
	if printName, e := run.Emplace(&PrintName{obj}); e != nil {
		err = e
	} else {
		err = run.ExecuteMatching(run, printName)
	}
	return
}

func printPluralName(run rt.Runtime, obj rt.Object) (err error) {
	if printName, e := run.Emplace(&PrintPluralName{obj}); e != nil {
		err = e
	} else {
		err = run.ExecuteMatching(run, printName)
	}
	return
}

func printSeveral(run rt.Runtime, obj rt.Object, cnt int) (err error) {
	if printName, e := run.Emplace(&PrintSeveral{obj, float64(cnt)}); e != nil {
		err = e
	} else {
		err = run.ExecuteMatching(run, printName)
	}
	return
}
