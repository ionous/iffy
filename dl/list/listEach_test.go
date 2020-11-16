package list_test

import (
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/kr/pretty"
)

type visitEach struct {
	visits *[]accum
}

func TestEach(t *testing.T) {
	// primary looping test:
	eachTest(t, []string{
		"Orange", "Lemon", "Mango",
	}, []accum{
		// what we expect to see from the index, first, last, and text values
		// when looping over the list of fruits
		{1, true, false, "Orange"},
		{2, false, false, "Lemon"},
		{3, false, true, "Mango"},
	}, 0 /*... and zero calls to else */)

	// never any middle ground
	eachTest(t, []string{
		"Orange", "Mango",
	}, []accum{
		{1, true, false, "Orange"},
		{2, false, true, "Mango"},
	}, 0 /*... and zero calls to else */)

	// first and last are both true
	eachTest(t, []string{
		"Lime",
	}, []accum{
		{1, true, true, "Lime"},
	}, 0 /*... and zero calls to else */)

	// looping over an empty list
	eachTest(t, nil, nil,
		1 /*... and a single call to else */)
}

func eachTest(t *testing.T, src []string, res []accum, otherwise int) {
	var out []string
	var visits []accum
	each := &list.Each{
		List: "src",
		With: "text",
		Go:   core.NewActivity(&visitEach{&visits}),
		Else: core.NewActivity(&Write{&out, T("x")}),
	}
	if e := each.Execute(&listTime{vals: values{"src": g.StringsOf(src)}}); e != nil {
		t.Fatal(src, e)
	} else if d := pretty.Diff(visits, res); len(d) > 0 || len(out) != otherwise {
		t.Fatal(src, out, d)
	}
}

func (v *visitEach) Execute(run rt.Runtime) (err error) {
	if i, e := index.GetNumber(run); e != nil {
		err = e
	} else if f, e := first.GetBool(run); e != nil {
		err = e
	} else if l, e := last.GetBool(run); e != nil {
		err = e
	} else if t, e := text.GetText(run); e != nil {
		err = e
	} else {
		(*v.visits) = append((*v.visits), accum{
			int(i), f, l, t,
		})
	}
	return
}

var (
	index rt.Variable = "index"
	first rt.Variable = "first"
	last  rt.Variable = "last"
	text  rt.Variable = "text"
)

type accum struct {
	index int
	first bool
	last  bool
	text  string
}
