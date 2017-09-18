package group

import (
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
	// a label means we want to treat the group as a block
	// including all pre and post test into that same block.
	block := len(g.Objects) > 0 && len(g.Label) > 0
	if !block {
		err = g.printGroup(run)
	} else {
		err = rt.WritersBlock(run, printer.Spanning(run.Writer()), func() error {
			return g.printGroup(run)
		})
	}
	return
}

func (g Collection) printGroup(run rt.Runtime) error {
	printGroup := run.Emplace(&PrintGroup{g.Label, g.Innumerable, g.ObjectGrouping, g.Objects})
	return run.ExecuteMatching(printGroup)
}
