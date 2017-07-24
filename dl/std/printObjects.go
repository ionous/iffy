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
		var buffer printer.Span
		if brackets {
			bracket := printer.Bracket{Writer: &buffer}
			run.PushWriter(&bracket)
		}
		// keep and separator outside of print so print could run by line as well.
		{
			var buffer printer.Span
			run.PushWriter(printer.AndSeparator(&buffer))

			if articles {
				err = printWithArticles(run, objs)
			} else {
				err = printWithoutArticles(run, objs)
			}

			run.PopWriter()
			run.Write(buffer.Bytes())
		}
		if brackets {
			run.PopWriter()
			run.Write(buffer.Bytes())
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
