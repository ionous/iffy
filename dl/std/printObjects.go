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
		if tersely {
			err = op.print(run, articles, objs)
		} else {
			err = rt.WritersBlock(run, printer.AndSeparator(run.Writer()), func() error {
				return op.print(run, articles, objs)
			})
		}
	}
	return
}

func (op *PrintObjects) print(run rt.Runtime, articles bool, objs rt.ObjectStream) (err error) {
	var cnt int
	if articles {
		cnt, err = printWithArticles(run, objs)
	} else {
		cnt, err = printWithoutArticles(run, objs)
	}
	if err == nil && cnt == 0 {
		err = op.Else.Execute(run)
	}
	return
}

func printWithArticles(run rt.Runtime, objs rt.ObjectStream) (ret int, err error) {
	// NoNames collects objects without names to print them in ordinal style.
	// This is needed for the edge case of a group of items which contains unnamed items.
	var nonames []rt.Object
	for ; objs.HasNext(); ret++ {
		if obj, e := objs.GetObject(); e != nil {
			err = e
			break
		} else {
			var name string
			if obj.GetValue("name", &name); e != nil {
				err = e
				break
			} else if unnamed := strings.Contains(name, "#"); unnamed {
				nonames = append(nonames, obj)
			} else if e := printArticle(run, "", obj); e != nil {
				err = e
				break
			}
		}
	}
	if err == nil {
		err = printSeveral(run, nonames)
	}
	return
}

func printWithoutArticles(run rt.Runtime, objs rt.ObjectStream) (ret int, err error) {
	var nonames []rt.Object
	for ; objs.HasNext(); ret++ {
		if obj, e := objs.GetObject(); e != nil {
			err = e
			break
		} else {
			var name string
			if obj.GetValue("name", &name); e != nil {
				err = e
				break
			} else if unnamed := strings.Contains(name, "#"); unnamed {
				nonames = append(nonames, obj)
			} else if e := rt.Determine(run, &PrintName{obj.Id()}); e != nil {
				err = e
				break
			}
		}
	}
	if err == nil {
		err = printSeveral(run, nonames)
	}
	return
}

// printSeveral objects without names in an ordinal style.
// This is needed for the edge case of a group of items which contains unnamed items.
func printSeveral(run rt.Runtime, objs []rt.Object) (err error) {
	if cnt := len(objs); cnt > 0 {
		err = rt.Determine(run, &PrintSeveral{objs[0].Id(), float64(cnt)})
	}
	return
}
