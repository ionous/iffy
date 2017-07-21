package std

import (
	"github.com/divan/num2words"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"io"
)

// sort the groups
type Sorted struct {
	*Group
	objs []rt.Object
}

func printGroups(run rt.Runtime, groups map[Group]ObjectList) (err error) {
	sorted := make([]Sorted, len(groups))
	for group, list := range groups {
		s := &(sorted[list.Order])
		s.Group = &group
		s.objs = list.Objects
	}
	// print each group
	for _, g := range sorted {
		var buffer printer.Span
		run.PushWriter(printer.AndSeparator(&buffer))

		if printGroupPattern, e := run.Emplace(&PrintGroup{g.Label, g.Innumerable, g.ObjectGrouping, g.objs}); e != nil {
			err = e
			break
		} else {
			if ran, e := run.ExecuteMatching(printGroupPattern); e != nil {
				err = e
				break
			} else if !ran {
				if e := defaultPrintGroup(run, g); e != nil {
					err = e
					break
				}
			}
		}

		run.PopWriter()
		// flush the group memebers to the output as a single unit
		if _, e := run.Write(buffer.Bytes()); e != nil {
			err = e
			break
		}
	}
	return
}

func defaultPrintGroup(run rt.Runtime, g Sorted) (err error) {
	var buffer printer.Span
	run.PushWriter(&buffer)
	//
	if len(g.Label) > 0 {
		n := len(g.objs)
		//
		if !g.Innumerable {
			if s := num2words.Convert(int(n)); len(s) > 0 {
				io.WriteString(&buffer, s)
			}
		}
		l := g.Label
		if n > 1 {
			l = run.Pluralize(l)
		}
		io.WriteString(&buffer, l)
	}
	//
	if g.ObjectGrouping != GroupWithoutObjects {
		bracket := printer.Bracket{Writer: &buffer}
		run.PushWriter(&bracket)
		//
		if g.ObjectGrouping == GroupWithArticles {
			err = printWithArticles(run, g.objs)
		} else {
			err = printWithoutArticles(run, g.objs)
		}
		run.PopWriter()
	}
	//
	run.PopWriter()
	_, e := run.Write(buffer.Bytes())
	return e
}

func printWithArticles(run rt.Runtime, objs []rt.Object) (err error) {
	for _, obj := range objs {
		if text, e := articleName(run, "", obj); e != nil {
			err = e
			break
		} else if _, e := io.WriteString(run, text); e != nil {
			err = e
			break
		}
	}
	return
}

func printWithoutArticles(run rt.Runtime, objs []rt.Object) (err error) {
	for _, obj := range objs {
		if e := printName(run, obj); e != nil {
			err = e
			break
		}
	}
	return
}
