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
		err = rt.WritersBlock(run, printer.AndSeparator(run.Writer()), func() (err error) {
			ungrouped := stream.NewObjectStream(stream.FromList(ungrouped))
			if _, e := printWithArticles(run, ungrouped); e != nil {
				err = e
			} else if e := groups.PrintGroups(run); e != nil {
				err = e
			}
			return
		})
	}
	return
}
