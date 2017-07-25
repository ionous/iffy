package std

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"strings"
)

type PrintObjects struct {
	Objects            rt.ObjListEval
	Articles, Brackets rt.BoolEval
}

func getBool(run rt.Runtime, eval rt.BoolEval) (ret bool, err error) {
	if eval != nil {
		if ok, e := eval.GetBool(run); e != nil {
			err = e
		} else {
			ret = ok
		}
	}
	return
}

func (op *PrintObjects) Execute(run rt.Runtime) (err error) {
	if objs, e := op.Objects.GetObjectStream(run); e != nil {
		err = e
	} else if articles, e := getBool(run, op.Articles); e != nil {
		err = e
	} else if brackets, e := getBool(run, op.Brackets); e != nil {
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
		err = printSeveral(run, n[0], cnt)
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
			break
		} else if unnamed := strings.Contains(kind.Name, "#"); unnamed {
			nonames = append(nonames, obj)
		} else if e := printArticle(run, "", obj); e != nil {
			err = e
			break
		}
	}
	if err == nil {
		err = nonames.Print(run)
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
			break
		} else if unnamed := strings.Contains(kind.Name, "#"); unnamed {
			nonames = append(nonames, obj)
		} else if e := printName(run, obj); e != nil {
			err = e
			break
		}
	}
	if err == nil {
		err = nonames.Print(run)
	}
	return
}
