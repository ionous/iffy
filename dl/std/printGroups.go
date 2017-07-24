package std

import (
	"github.com/divan/num2words"
	"github.com/ionous/iffy/dl/std/group"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"io"
)

func printGroups(run rt.Runtime, groups []group.Collection) (err error) {
	for _, g := range groups {
		if len(g.Objects) > 0 {
			var buffer printer.Span
			run.PushWriter(printer.AndSeparator(&buffer))

			if printGroupPattern, e := run.Emplace(&group.PrintGroup{g.Label, g.Innumerable, g.ObjectGrouping, g.Objects}); e != nil {
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
	}
	return
}

// FIX? move into a pattern
func defaultPrintGroup(run rt.Runtime, g group.Collection) (err error) {
	var buffer printer.Span
	run.PushWriter(&buffer)
	//
	if len(g.Label) > 0 {
		n := len(g.Objects)
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
	if g.ObjectGrouping != group.WithoutObjects {
		bracket := printer.Bracket{Writer: &buffer}
		run.PushWriter(&bracket)
		//
		if g.ObjectGrouping == group.WithArticles {
			err = printWithArticles(run, g.Objects)
		} else {
			err = printWithoutArticles(run, g.Objects)
		}
		run.PopWriter()
	}
	//
	run.PopWriter()
	_, e := run.Write(buffer.Bytes())
	return e
}

// FIX? move into a pattern
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

// FIX? move into a pattern
func printWithoutArticles(run rt.Runtime, objs []rt.Object) (err error) {
	for _, obj := range objs {
		if e := printName(run, obj); e != nil {
			err = e
			break
		}
	}
	return
}
