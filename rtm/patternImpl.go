package rtm

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

func (run *Rtm) GetBoolMatching(data rt.Object) (ret bool, err error) {
	rt.ScopeBlock(run, data, func() {
		rtype := data.Type()
		if ps, ok := run.rules.Bools[rtype]; !ok {
			err = notFound(rtype)
		} else {
			ret, err = ps.GetBool(run)
		}
	})
	return
}
func (run *Rtm) GetNumMatching(data rt.Object) (ret float64, err error) {
	rt.ScopeBlock(run, data, func() {
		rtype := data.Type()
		if ps, ok := run.rules.Numbers[rtype]; !ok {
			err = notFound(rtype)
		} else {
			ret, err = ps.GetNumber(run)
		}
	})
	return
}
func (run *Rtm) GetTextMatching(data rt.Object) (ret string, err error) {
	rt.ScopeBlock(run, data, func() {
		rtype := data.Type()
		if ps, ok := run.rules.Text[rtype]; !ok {
			err = notFound(rtype)
		} else {
			ret, err = ps.GetText(run)
		}
	})
	return
}
func (run *Rtm) GetObjectMatching(data rt.Object) (ret rt.Object, err error) {
	rt.ScopeBlock(run, data, func() {
		rtype := data.Type()
		if ps, ok := run.rules.Objects[rtype]; !ok {
			err = notFound(rtype)
		} else {
			ret, err = ps.GetObject(run)
		}
	})
	return
}
func (run *Rtm) GetNumStreamMatching(data rt.Object) (ret rt.NumberStream, err error) {
	rt.ScopeBlock(run, data, func() {
		rtype := data.Type()
		if ps, ok := run.rules.NumLists[rtype]; !ok {
			err = notFound(rtype)
		} else {
			ret, err = ps.GetNumberStream(run)
		}
	})
	return
}
func (run *Rtm) GetTextStreamMatching(data rt.Object) (ret rt.TextStream, err error) {
	rt.ScopeBlock(run, data, func() {
		rtype := data.Type()
		if ps, ok := run.rules.TextLists[rtype]; !ok {
			err = notFound(rtype)
		} else {
			ret, err = ps.GetTextStream(run)
		}
	})
	return
}
func (run *Rtm) GetObjStreamMatching(data rt.Object) (ret rt.ObjectStream, err error) {
	rt.ScopeBlock(run, data, func() {
		rtype := data.Type()
		if ps, ok := run.rules.ObjLists[rtype]; !ok {
			err = notFound(rtype)
		} else {
			ret, err = ps.GetObjectStream(run)
		}
	})
	return
}
func (run *Rtm) ExecuteMatching(data rt.Object) (err error) {
	rt.ScopeBlock(run, data, func() {
		rtype := data.Type()

		// NOTE: if we need to differentiate between "ran" and "not found",
		// "didnt run" should probably become an error code.
		if ps, ok := run.rules.Executes[rtype]; ok {
			_, err = ps.Execute(run)
		}
	})
	return
}

func notFound(rtype r.Type) error {
	return errutil.New("no rules found for", rtype)
}
