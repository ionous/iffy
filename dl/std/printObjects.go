package std

import (
	"github.com/divan/num2words"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"io"
	"strings"
)

type PrintObjects struct {
	Objects            rt.ObjListEval
	Articles, Brackets rt.BoolEval
}

func (op *PrintObjects) Execute(run rt.Runtime) (err error) {
	if objs, e := op.Objects.GetObjectStream(run); e != nil {
		err = e
	} else if articles, e := op.Articles.GetBool(run); e != nil {
		err = e
	} else if brackets, e := op.Brackets.GetBool(run); e != nil {
		err = e
	} else {
		if brackets {
			bracket := printer.Bracket{Writer: run}
			run = rt.Writer(run, &bracket)
			defer bracket.Flush()
		}
		// keep and separator outside of print so print could run by line as well.
		sep := printer.AndSeparator(run)
		run = rt.Writer(run, sep)
		defer sep.Flush()
		//
		if articles {
			err = printWithArticles(run, objs)
		} else {
			err = printWithoutArticles(run, objs)
		}
	}
	return
}

// NoNames collects objects without names to print them in ordinal style.
// This is needed for the edge case of a group of items which contains unnamed items.
type NoNames []rt.Object

func (n NoNames) Print(run rt.Runtime) (err error) {
	if cnt := len(n); cnt > 0 {
		// FIX: "print plural name" --> maybe the count should go in the pattern
		span := printer.Spanner{Writer: run}
		run = rt.Writer(run, &span)
		num := num2words.Convert(cnt)
		if _, e := io.WriteString(run, num); e != nil {
			err = e
		} else {
			err = printPluralName(run, n[0])
		}
		if err == nil {
			err = span.Flush()
		}
	}
	return
}

func printWithArticles(run rt.Runtime, objs rt.ObjectStream) (err error) {
	var nonames NoNames
	for objs.HasNext() {
		if obj, e := objs.GetNext(); e != nil {
			err = e
			break
		} else if kind, e := kindOf(run, obj); e != nil {
			err = e
		} else if unnamed := strings.Contains(kind.Name, "#"); unnamed {
			nonames = append(nonames, obj)
		} else if text, e := articleName(run, "", obj); e != nil {
			err = e
			break
		} else if _, e := io.WriteString(run, text); e != nil {
			err = e
			break
		}
	}
	if err == nil {
		nonames.Print(run)
	}
	return
}

func printWithoutArticles(run rt.Runtime, objs rt.ObjectStream) (err error) {
	var nonames NoNames
	for objs.HasNext() {
		if obj, e := objs.GetNext(); e != nil {
			err = e
			break
		} else if kind, e := kindOf(run, obj); e != nil {
			err = e
		} else if unnamed := strings.Contains(kind.Name, "#"); unnamed {
			nonames = append(nonames, obj)
		} else if text, e := getName(run, obj); e != nil {
			err = e
			break
		} else if _, e := io.WriteString(run, text); e != nil {
			err = e
			break
		}
	}
	if err == nil {
		nonames.Print(run)
	}
	return
}
