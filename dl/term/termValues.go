package term

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

type Number struct {
	Name string        // parameter name
	Init rt.NumberEval // default value
}
type Bool struct {
	Name string
	Init rt.BoolEval
}
type Text struct {
	Name string
	Init rt.TextEval
}
type Record struct {
	Name, Kind string
}
type NumList struct {
	Name string
	Init rt.NumListEval
}
type TextList struct {
	Name string
	Init rt.TextListEval
}
type RecordList struct {
	Name, Kind string
}

func (n *Number) String() string {
	return n.Name
}

func (n *Number) Prepare(run rt.Runtime, p *Terms) (err error) {
	if v, e := rt.GetOptionalNumber(run, n.Init, 0); e != nil {
		err = e
	} else if t, e := g.ValueOf(v); e != nil {
		err = e
	} else {
		p.AddTerm(n.Name, t)
	}
	return
}

func (n *Bool) String() string {
	return n.Name
}

func (n *Bool) Prepare(run rt.Runtime, p *Terms) (err error) {
	if v, e := rt.GetOptionalBool(run, n.Init, false); e != nil {
		err = e
	} else if t, e := g.ValueOf(v); e != nil {
		err = e
	} else {
		p.AddTerm(n.Name, t)
	}
	return
}

func (n *Text) String() string {
	return n.Name
}

func (n *Text) Prepare(run rt.Runtime, p *Terms) (err error) {
	if v, e := rt.GetOptionalText(run, n.Init, ""); e != nil {
		err = e
	} else if t, e := g.ValueOf(v); e != nil {
		err = e
	} else {
		p.AddTerm(n.Name, t)
	}
	return
}

func (n *Record) String() string {
	return n.Name
}

func (n *Record) Prepare(run rt.Runtime, p *Terms) (err error) {
	if k, e := run.GetKindByName(n.Kind); e != nil {
		err = e
	} else if t, e := g.ValueOf(k.NewRecord()); e != nil {
		err = e
	} else {
		p.AddTerm(n.Name, t)
	}
	return
}

func (n *NumList) String() string {
	return n.Name
}

func (n *NumList) Prepare(run rt.Runtime, p *Terms) (err error) {
	if vs, e := rt.GetOptionalNumbers(run, n.Init, nil); e != nil {
		err = e
	} else if t, e := g.ValueOf(vs); e != nil {
		err = e
	} else {
		p.AddTerm(n.Name, t)
	}
	return
}

func (n *TextList) String() string {
	return n.Name
}

func (n *TextList) Prepare(run rt.Runtime, p *Terms) (err error) {
	if vs, e := rt.GetOptionalTexts(run, n.Init, nil); e != nil {
		err = e
	} else if t, e := g.ValueOf(vs); e != nil {
		err = e
	} else {
		p.AddTerm(n.Name, t)
	}
	return
}

func (n *RecordList) String() string {
	return n.Name
}

func (n *RecordList) Prepare(run rt.Runtime, p *Terms) (err error) {
	if k, e := run.GetKindByName(n.Kind); e != nil {
		err = e
	} else if t, e := g.ValueFrom([]*g.Record{}, affine.RecordList, k.Name()); e != nil {
		err = e
	} else {
		p.AddTerm(n.Name, t)
	}
	return
}
