package std

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"io"
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

func printWithArticles(run rt.Runtime, objs rt.ObjectStream) (err error) {
	for objs.HasNext() {
		if obj, e := objs.GetNext(); e != nil {
			err = e
			break
		} else if text, e := articleName(run, "", obj); e != nil {
			err = e
			break
		} else if _, e := io.WriteString(run, text); e != nil {
			err = e
			break
		}
	}
	return
}

func printWithoutArticles(run rt.Runtime, objs rt.ObjectStream) (err error) {
	for objs.HasNext() {
		if obj, e := objs.GetNext(); e != nil {
			err = e
			break
		} else if e := printName(run, obj); e != nil {
			err = e
			break
		}
	}
	return
}
