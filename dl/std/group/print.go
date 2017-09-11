package group

import (
	"github.com/ionous/iffy/ref/obj"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
)

func (c Collections) PrintGroups(run rt.Runtime) (err error) {
	for _, g := range c {
		if e := g.PrintGroup(run); e != nil {
			err = e
			break
		}
	}
	return
}

func (g Collection) PrintGroup(run rt.Runtime) (err error) {
	if len(g.Objects) > 0 {
		// a label means we want to treat the group as a block
		// including all pre and post test into that same block.
		if len(g.Label) > 0 {
			span := printer.Spanner{Writer: run}
			run = rt.Writer(run, &span)
			defer span.Flush()
		}
		printGroup := obj.Emplace(&PrintGroup{g.Label, g.Innumerable, g.ObjectGrouping, g.Objects})
		err = run.ExecuteMatching(run, printGroup)
	}

	return
}
