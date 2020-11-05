package term

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
)

type Number struct {
	Name string
	Init rt.NumberEval
}
type Bool struct {
	Name string
	Init rt.BoolEval
}
type Text struct {
	Name string
	Init rt.TextEval
}
type Object struct {
	Name, Kind string
	Init       rt.TextEval
}
type NumList struct {
	Name string
	Init rt.NumListEval
}
type TextList struct {
	Name string
	Init rt.TextListEval
}

// type RecordList struct {
// 	Name, Kind string
// 	// possibly with an initial size generating a zero list.
// }

func (n *Number) String() string {
	return n.Name
}

func (n *Number) Prepare(run rt.Runtime, p *Terms) (err error) {
	if v, e := rt.GetOptionalNumber(run, n.Init, 0); e != nil {
		err = e
	} else {
		p.AddTerm(n.Name, generic.NewFloat(v))
	}
	return
}

func (n *Bool) String() string {
	return n.Name
}

func (n *Bool) Prepare(run rt.Runtime, p *Terms) (err error) {
	if v, e := rt.GetOptionalBool(run, n.Init, false); e != nil {
		err = e
	} else {
		p.AddTerm(n.Name, generic.NewBool(v))
	}
	return
}

func (n *Text) String() string {
	return n.Name
}

func (n *Text) Prepare(run rt.Runtime, p *Terms) (err error) {
	if v, e := rt.GetOptionalText(run, n.Init, ""); e != nil {
		err = e
	} else {
		p.AddTerm(n.Name, generic.NewString(v))
	}
	return
}

func (n *Object) String() string {
	return n.Name
}

func (n *Object) Prepare(run rt.Runtime, p *Terms) (err error) {
	if v, e := rt.GetOptionalText(run, n.Init, ""); e != nil {
		err = e
	} else {
		p.AddTerm(n.Name, generic.NewString(v))
	}
	return
}

func (n *NumList) String() string {
	return n.Name
}

func (n *NumList) Prepare(run rt.Runtime, p *Terms) (err error) {
	if vs, e := rt.GetOptionalNumbers(run, n.Init, nil); e != nil {
		err = e
	} else {
		p.AddTerm(n.Name, generic.NewFloatSlice(vs))
	}
	return
}

func (n *TextList) String() string {
	return n.Name
}

func (n *TextList) Prepare(run rt.Runtime, p *Terms) (err error) {
	if vs, e := rt.GetOptionalTexts(run, n.Init, nil); e != nil {
		err = e
	} else {
		p.AddTerm(n.Name, generic.NewStringSlice(vs))
	}
	return
}

// func (n *RecordList) String() string {
// 	return n.Name
// }

// func (n *RecordList) Prepare(run rt.Runtime, p *Terms) (err error) {
// 	p.AddTerm(n.Name, generic.NewNewRecordSlice(n.Kind))
// 	return
// }
