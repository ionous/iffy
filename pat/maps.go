package pat

import (
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/class"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
)

// FIX? consider rewriting using reflect so that all patterns are stored together and if the interface doesnt match the expectation it errors.
type Bools map[ident.Id]BoolRules
type Numbers map[ident.Id]NumberRules
type Text map[ident.Id]TextRules
type Objects map[ident.Id]ObjectRules
type NumLists map[ident.Id]NumListRules
type TextLists map[ident.Id]TextListRules
type ObjLists map[ident.Id]ObjListRules
type Executes map[ident.Id]ExecuteRules

func setupScope(run rt.Runtime, data rt.Object, cb func(ident.Id)) {
	id := class.Id(data.Type())
	run.PushScope(scope.AtFinder(data))
	cb(id)
	run.PopScope()
}

func (m Bools) AddRule(id ident.Id, f Filters, k rt.BoolEval) {
	m[id] = append(m[id], BoolRule{f, k})
}
func (m Numbers) AddRule(id ident.Id, f Filters, k rt.NumberEval) {
	m[id] = append(m[id], NumberRule{f, k})
}
func (m Text) AddRule(id ident.Id, f Filters, k rt.TextEval) {
	m[id] = append(m[id], TextRule{f, k})
}
func (m Objects) AddRule(id ident.Id, f Filters, k rt.ObjectEval) {
	m[id] = append(m[id], ObjectRule{f, k})
}
func (m NumLists) AddRule(id ident.Id, f Filters, k rt.NumListEval) {
	m[id] = append(m[id], NumListRule{f, k})
}
func (m TextLists) AddRule(id ident.Id, f Filters, k rt.TextListEval) {
	m[id] = append(m[id], TextListRule{f, k})
}
func (m ObjLists) AddRule(id ident.Id, f Filters, k rt.ObjListEval) {
	m[id] = append(m[id], ObjListRule{f, k})
}
func (m Executes) AddRule(id ident.Id, f Filters, k rt.Execute, flags Flags) {
	m[id] = append(m[id], ExecuteRule{f, k, flags})
}

func (m Bools) GetBoolMatching(run rt.Runtime, data rt.Object) (ret bool, err error) {
	setupScope(run, data, func(id ident.Id) {
		if ps, ok := m[id]; !ok {
			err = NotFound(id.Name)
		} else {
			ret, err = ps.GetBool(run)
		}
	})
	return
}
func (m Numbers) GetNumMatching(run rt.Runtime, data rt.Object) (ret float64, err error) {
	setupScope(run, data, func(id ident.Id) {
		if ps, ok := m[id]; !ok {
			err = NotFound(id.Name)
		} else {
			ret, err = ps.GetNumber(run)
		}
	})
	return
}
func (m Text) GetTextMatching(run rt.Runtime, data rt.Object) (ret string, err error) {
	setupScope(run, data, func(id ident.Id) {
		if ps, ok := m[id]; !ok {
			err = NotFound(id.Name)
		} else {
			ret, err = ps.GetText(run)
		}
	})
	return
}

func (m Objects) GetObjectMatching(run rt.Runtime, data rt.Object) (ret rt.Object, err error) {
	setupScope(run, data, func(id ident.Id) {
		if ps, ok := m[id]; !ok {
			err = NotFound(id.Name)
		} else {
			ret, err = ps.GetObject(run)
		}
	})
	return
}
func (m NumLists) GetNumStreamMatching(run rt.Runtime, data rt.Object) (ret rt.NumberStream, err error) {
	setupScope(run, data, func(id ident.Id) {
		if ps, ok := m[id]; !ok {
			err = NotFound(id.Name)
		} else {
			ret, err = ps.GetNumberStream(run)
		}
	})
	return
}

func (m TextLists) GetTextStreamMatching(run rt.Runtime, data rt.Object) (ret rt.TextStream, err error) {
	setupScope(run, data, func(id ident.Id) {
		if ps, ok := m[id]; !ok {
			err = NotFound(id.Name)
		} else {
			ret, err = ps.GetTextStream(run)
		}
	})
	return
}

func (m ObjLists) GetObjStreamMatching(run rt.Runtime, data rt.Object) (ret rt.ObjectStream, err error) {
	setupScope(run, data, func(id ident.Id) {
		if ps, ok := m[id]; !ok {
			err = NotFound(id.Name)
		} else {
			ret, err = ps.GetObjectStream(run)
		}
	})
	return
}

func (m Executes) ExecuteMatching(run rt.Runtime, data rt.Object) (err error) {
	setupScope(run, data, func(id ident.Id) {
		// NOTE: if we need to differentiate between "ran" and "not found",
		// "didnt run" should become an error code.
		if ps, ok := m[id]; ok {
			_, err = ps.Execute(run)
		}
	})
	return
}
