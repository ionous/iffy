package pat

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
)

// Flags controlling how individual list rules ( which each respond with multiple results ) work together.
type Flags int

//go:generate stringer -type=Flags
const (
	Infix Flags = iota
	Prefix
	Postfix
)

// BoolRule responds with a true/false result when its filters are satisfied.
// It works in conjunction with BoolRules.
type BoolRule struct {
	Filters Filters
	rt.BoolEval
}

// NumberRule responds with a single number when its filters are satisfied.
// It works in conjunction with NumberRules.
type NumberRule struct {
	Filters Filters
	rt.NumberEval
}

// TextRule responds with a bit of text when its filters are satisfied.
// It works in conjunction with TextRules.
type TextRule struct {
	Filters Filters
	rt.TextEval
}

// ObjectRule responds with an object when its filters are satisfied.
// It works in conjunction with ObjectRules.
type ObjectRule struct {
	Filters Filters
	rt.ObjectEval
}

// ListRule for any rule which can respond with multiple results.
type ListRule struct {
	Filters Filters
	Flags
}

// NumListRule responds with a stream of numbers when its filters are satisfied.
// It works in conjunction with NumListRules.
type NumListRule struct {
	ListRule
	rt.NumListEval
}

// TextListRule responds with a stream of text when its filters are satisfied.
// It works in conjunction with TextListRules.
type TextListRule struct {
	ListRule
	rt.TextListEval
}

// ObjListRule responds with a stream of objects when its filters are satisfied.
// It works in conjunction with ObjListRules.
type ObjListRule struct {
	ListRule
	rt.ObjListEval
}

// ExecuteListRule triggers a series of statements when its filters are satisfied.
// It works in conjunction with ExecuteRules.
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
