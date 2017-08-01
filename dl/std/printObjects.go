package std

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/optional"
	"github.com/ionous/iffy/rt/printer"
	"io"
	"strings"
)

type PrintObjects struct {
	Objects           rt.ObjListEval
	Header            rt.TextEval
	Articles, Tersely rt.BoolEval
	Else              rt.ExecuteList
}

func (op *PrintObjects) Execute(run rt.Runtime) (err error) {
	if objs, e := op.Objects.GetObjectStream(run); e != nil {
		err = e
	} else if articles, e := optional.Bool(run, op.Articles); e != nil {
		err = e
	} else if tersely, e := optional.Bool(run, op.Tersely); e != nil {
		err = e
	} else if header, e := optional.Text(run, op.Header); e != nil {
		err = e
	} else {
		if len(header) > 0 {
			io.WriteString(run, header)
		}
		// control and separator so print could run by line as well.
		if !tersely {
			sep := printer.AndSeparator(run)
			run = rt.Writer(run, sep)
			defer sep.Flush()
		}
		//
		var cnt int
		if articles {
			cnt, err = printWithArticles(run, objs)
		} else {
			cnt, err = printWithoutArticles(run, objs)
		}
		if err == nil && cnt == 0 {
			err = op.Else.Execute(run)
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

func printWithArticles(run rt.Runtime, objs rt.ObjectStream) (ret int, err error) {
	var nonames NoNames
	for ; objs.HasNext(); ret++ {
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

func printWithoutArticles(run rt.Runtime, objs rt.ObjectStream) (ret int, err error) {
	var nonames NoNames
	for ; objs.HasNext(); ret++ {
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
