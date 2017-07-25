package group

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
)

func (c Collections) Print(run rt.Runtime) (err error) {
	for _, g := range c {
		if e := g.Print(run); e != nil {
			err = e
			break
		}
	}
	return
}

func (g Collection) Print(run rt.Runtime) (err error) {
	if len(g.Objects) > 0 {
		// a label means we want to treat the group as a block
		// including all pre and post test into that same block.
		if len(g.Label) > 0 {
			span := printer.Spanner{Writer: run}
			run = rt.Writer(run, &span)
			defer span.Flush()
		}

		if printGroup, e := run.Emplace(&PrintGroup{g.Label, g.Innumerable, g.ObjectGrouping, g.Objects}); e != nil {
			err = e
		} else if e := run.ExecuteMatching(run, printGroup); e != nil {
			err = e
		}
	}
	return
}
