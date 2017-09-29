package rtm

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

func (run *Rtm) GetBoolMatching(data rt.Object) (ret bool, err error) {
	rules := run.Rules
	rt.ScopeBlock(run, data, func() {
		rtype := data.Type()
		if ps, ok := rules.Bools[rtype]; !ok {
			err = notFound(rtype)
		} else {
			ret, err = ps.GetBool(run)
		}
	})
	return
}
func (run *Rtm) GetNumMatching(data rt.Object) (ret float64, err error) {
	rules := run.Rules
	rt.ScopeBlock(run, data, func() {
		rtype := data.Type()
		if ps, ok := rules.Numbers[rtype]; !ok {
			err = notFound(rtype)
		} else {
			ret, err = ps.GetNumber(run)
		}
	})
	return
}
func (run *Rtm) GetTextMatching(data rt.Object) (ret string, err error) {
	rules := run.Rules
	rt.ScopeBlock(run, data, func() {
		rtype := data.Type()
		if ps, ok := rules.TextPatterns[rtype]; !ok {
			err = notFound(rtype)
		} else {
			ret, err = ps.GetText(run)
		}
	})
	return
}
func (run *Rtm) GetObjectMatching(data rt.Object) (ret rt.Object, err error) {
	rules := run.Rules
	rt.ScopeBlock(run, data, func() {
		rtype := data.Type()
		if ps, ok := rules.Objects[rtype]; !ok {
			err = notFound(rtype)
		} else {
			ret, err = ps.GetObject(run)
		}
	})
	return
}
func (run *Rtm) GetNumStreamMatching(data rt.Object) (ret rt.NumberStream, err error) {
	rules := run.Rules
	rt.ScopeBlock(run, data, func() {
		rtype := data.Type()
		if ps, ok := rules.NumLists[rtype]; !ok {
			err = notFound(rtype)
		} else {
			ret, err = ps.GetNumberStream(run)
		}
	})
	return
}
func (run *Rtm) GetTextStreamMatching(data rt.Object) (ret rt.TextStream, err error) {
	rules := run.Rules
	rt.ScopeBlock(run, data, func() {
		rtype := data.Type()
		if ps, ok := rules.TextLists[rtype]; !ok {
			err = notFound(rtype)
		} else {
			ret, err = ps.GetTextStream(run)
		}
	})
	return
}
func (run *Rtm) GetObjStreamMatching(data rt.Object) (ret rt.ObjectStream, err error) {
	rules := run.Rules
	rt.ScopeBlock(run, data, func() {
		rtype := data.Type()
		if ps, ok := rules.ObjLists[rtype]; !ok {
			err = notFound(rtype)
		} else {
			ret, err = ps.GetObjectStream(run)
		}
	})
	return
}
func (run *Rtm) ExecuteMatching(data rt.Object) (err error) {
	rules := run.Rules
	rt.ScopeBlock(run, data, func() {
		rtype := data.Type()
		if ps, ok := rules.Executes[rtype]; ok {
			// NOTE: if we need to differentiate between "ran" and "not found",
			// "didnt run" should probably become an error code.
			_, err = ps.Execute(run)
		}
	})
	return
}

func notFound(rtype r.Type) error {
	return errutil.New("no rules found for", rtype)
}
