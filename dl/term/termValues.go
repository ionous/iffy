package term

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/rt"
)

type Number struct {
	Name string // parameter name
}
type Bool struct {
	Name string
}
type Text struct {
	Name string
}
type Record struct {
	Name, Kind string
}
type NumList struct {
	Name string
}
type TextList struct {
	Name string
}
type RecordList struct {
	Name, Kind string
}

func (n *Number) String() string {
	return n.Name
}

func (n *Number) Prepare(run rt.Runtime, p *Terms) (err error) {
	p.AddTerm(n.Name, affine.Number, "")
	return
}

func (n *Bool) String() string {
	return n.Name
}

func (n *Bool) Prepare(run rt.Runtime, p *Terms) (err error) {
	p.AddTerm(n.Name, affine.Bool, "")
	return
}

func (n *Text) String() string {
	return n.Name
}

func (n *Text) Prepare(run rt.Runtime, p *Terms) (err error) {
	p.AddTerm(n.Name, affine.Text, "")
	return
}

func (n *Record) String() string {
	return n.Name
}

func (n *Record) Prepare(run rt.Runtime, p *Terms) (err error) {
	if k, e := run.GetKindByName(n.Kind); e != nil {
		err = e
	} else {
		p.AddTerm(n.Name, affine.Record, k.Name())
	}
	return
}

func (n *NumList) String() string {
	return n.Name
}

func (n *NumList) Prepare(run rt.Runtime, p *Terms) (err error) {
	p.AddTerm(n.Name, affine.NumList, "")
	return
}

func (n *TextList) String() string {
	return n.Name
}

func (n *TextList) Prepare(run rt.Runtime, p *Terms) (err error) {
	p.AddTerm(n.Name, affine.TextList, "")
	return
}

func (n *RecordList) String() string {
	return n.Name
}

func (n *RecordList) Prepare(run rt.Runtime, p *Terms) (err error) {
	if k, e := run.GetKindByName(n.Kind); e != nil {
		err = e
	} else {
		p.AddTerm(n.Name, affine.RecordList, k.Name())
	}
	return
}
