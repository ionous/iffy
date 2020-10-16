package term

import (
	"github.com/ionous/iffy/affine"
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

func (n *Number) String() string {
	return n.Name
}

func (n *Number) Prepare(run rt.Runtime, p *Terms) (err error) {
	if v, e := rt.GetOptionalNumber(run, n.Init, 0); e != nil {
		err = e
	} else {
		p.addTerm(n.Name, affine.Number, &generic.Float{Value: v})
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
		p.addTerm(n.Name, affine.Bool, &generic.Bool{Value: v})
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
		p.addTerm(n.Name, affine.Text, &generic.String{Value: v})
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
		p.addTerm(n.Name, affine.Text, &generic.String{Value: v})
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
		p.addTerm(n.Name, affine.NumList, &generic.FloatSlice{Value: vs})
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
		p.addTerm(n.Name, affine.TextList, &generic.StringSlice{Value: vs})
	}
	return
}
