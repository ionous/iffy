package list_test

import (
	"testing"

	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/safe"
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
		List: V("Source"),
		As:   &list.AsTxt{N("text")},
		Do:   core.MakeActivity(&visitEach{&visits}),
		Else: core.NewActivity(&Write{&out, T("x")}),
	}
	if lt, _, e := newListTime(src, nil); e != nil {
		t.Fatal(e)
	} else if e := each.Execute(lt); e != nil {
		t.Fatal(src, e)
	} else if d := pretty.Diff(visits, res); len(d) > 0 || len(out) != otherwise {
		t.Fatal(src, out, d)
	}
}

func (v *visitEach) Execute(run rt.Runtime) (err error) {
	if i, e := safe.Variable(run, "index", affine.Number); e != nil {
		err = e
	} else if f, e := safe.Variable(run, "first", affine.Bool); e != nil {
		err = e
	} else if l, e := safe.Variable(run, "last", affine.Bool); e != nil {
		err = e
	} else if t, e := safe.Variable(run, "text", affine.Text); e != nil {
		err = e
	} else {
		(*v.visits) = append((*v.visits), accum{
			i.Int(), f.Bool(), l.Bool(), t.String(),
		})
	}
	return
}

type accum struct {
	index int
	first bool
	last  bool
	text  string
}
