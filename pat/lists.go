package pat

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
)

type Flags int

//go:generate stringer -type=Flags
const (
	Infix Flags = iota
	Prefix
	Postfix
)

// BoolRule holds the distallation of a BoolRule.
type BoolRule struct {
	Filters Filters
	rt.BoolEval
}

// NumberRule holds the distallation of a NumberRule.
type NumberRule struct {
	Filters Filters
	rt.NumberEval
}

// TextRule holds the distallation of a TextRule.
type TextRule struct {
	Filters Filters
	rt.TextEval
}

// ObjectRule holds the distallation of a ObjectRule.
type ObjectRule struct {
	Filters Filters
	rt.ObjectEval
}

// ListRule base for all rules which return streams of data.
type ListRule struct {
	Filters Filters
	Flags
}

// NumListRule holds the distallation of a NumListRule.
type NumListRule struct {
	ListRule
	rt.NumListEval
}
type TextListRule struct {
	ListRule
	rt.TextListEval
}
type ObjListRule struct {
	ListRule
	rt.ObjListEval
}

// ExecuteRule holds the distallation of a RunRule.
type ExecuteRule struct {
	ListRule
	rt.Execute
}

type BoolRules []BoolRule
type NumberRules []NumberRule
type TextRules []TextRule
type ObjectRules []ObjectRule
type NumListRules []NumListRule
type TextListRules []TextListRule
type ObjListRules []ObjListRule
type ExecuteRules []ExecuteRule

func (ps BoolRules) GetBool(run rt.Runtime) (ret bool, err error) {
	for i, cnt := 0, len(ps); i < cnt; i++ {
		p := ps[cnt-i-1]
		if matched, e := p.Filters.GetBool(run); e != nil {
			err = e
			break
		} else if matched {
			ret, err = p.BoolEval.GetBool(run)
			break
		}
	}
	return
}
func (ps NumberRules) GetNumber(run rt.Runtime) (ret float64, err error) {
	for i, cnt := 0, len(ps); i < cnt; i++ {
		p := ps[cnt-i-1]
		if matched, e := p.Filters.GetBool(run); e != nil {
			err = e
			break
		} else if matched {
			ret, err = p.NumberEval.GetNumber(run)
			break
		}
	}
	return
}
func (ps TextRules) GetText(run rt.Runtime) (ret string, err error) {
	for i, cnt := 0, len(ps); i < cnt; i++ {
		p := ps[cnt-i-1]
		if matched, e := p.Filters.GetBool(run); e != nil {
			err = e
			break
		} else if matched {
			ret, err = p.TextEval.GetText(run)
			break
		}
	}
	return
}
func (ps ObjectRules) GetObject(run rt.Runtime) (ret rt.Object, err error) {
	for i, cnt := 0, len(ps); i < cnt; i++ {
		p := ps[cnt-i-1]
		if matched, e := p.Filters.GetBool(run); e != nil {
			err = e
			break
		} else if matched {
			ret, err = p.ObjectEval.GetObject(run)
			break
		}
	}
	return
}
func (r ListRule) Applies(run rt.Runtime) (ret Flags, err error) {
	if ok, e := r.Filters.GetBool(run); e != nil {
		err = e
	} else if !ok {
		ret = -1
	} else {
		ret = r.Flags
	}
	return
}
func (ps NumListRules) GetNumberStream(run rt.Runtime) (ret rt.NumberStream, err error) {
	if q, e := splitQuery(run, ps); e != nil {
		err = e
	} else {
		q := adaptNumbers(run, q)
		ret = stream.NewNumberStream(q.Iterate())
	}
	return
}
func (ps TextListRules) GetTextStream(run rt.Runtime) (ret rt.TextStream, err error) {
	if q, e := splitQuery(run, ps); e != nil {
		err = e
	} else {
		q := adaptText(run, q)
		ret = stream.NewTextStream(q.Iterate())
	}
	return
}
func (ps ObjListRules) GetObjectStream(run rt.Runtime) (ret rt.ObjectStream, err error) {
	if q, e := splitQuery(run, ps); e != nil {
		err = e
	} else {
		q := adaptObjects(run, q)
		ret = stream.NewObjectStream(q.Iterate())
	}
	return
}
func (ps ExecuteRules) Execute(run rt.Runtime) (ret bool, err error) {
	if q, e := splitQuery(run, ps); e != nil {
		err = e
	} else {
		next := q.Iterate()
		for item, ok := next(); ok; item, ok = next() {
			if e := item.(ExecuteRule).Execute.Execute(run); e != nil {
				err = e
				break
			}
			ret = true // any executed
		}
	}
	return
}
