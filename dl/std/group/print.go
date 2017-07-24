package group

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
)

func (c Collections) Print(run rt.Runtime) (err error) {
	for _, g := range c {
		if len(g.Objects) > 0 {
			var buffer printer.Span
			run.PushWriter(printer.AndSeparator(&buffer))
			//
			if printGroup, e := run.Emplace(&PrintGroup{g.Label, g.Innumerable, g.ObjectGrouping, g.Objects}); e != nil {
				err = e
				break
			} else if _, e := run.ExecuteMatching(printGroup); e != nil {
				err = e
				break
			}
			//
			run.PopWriter()
			// flush the group memebers to the output as a single unit
			if _, e := run.Write(buffer.Bytes()); e != nil {
				err = e
				break
			}
		}
	}
	return
}
