package pat

import (
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
)

// FIX? consider rewriting using reflect so that all patterns are stored together and if the interface doesnt match the expectation it errors.
type BoolMap map[string]BoolPatterns
type NumberMap map[string]NumberPatterns
type TextMap map[string]TextPatterns
type ObjectMap map[string]ObjectPatterns
type NumListMap map[string]NumListPatterns
type TextListMap map[string]TextListPatterns
type ObjListMap map[string]ObjListPatterns

func (m BoolMap) GetBoolMatching(run rt.Runtime, name string, data rt.Object) (ret bool, err error) {
	id := id.MakeId(name)
	if ps, ok := m[id]; !ok {
		err = NotFound
	} else {
		run.PushScope(scope.AtFinder(data))
		ret, err = ps.GetBool(run)
		run.PopScope()
	}
	return
}

func (m NumberMap) GetNumMatching(run rt.Runtime, name string, data rt.Object) (ret float64, err error) {
	id := id.MakeId(name)
	if ps, ok := m[id]; !ok {
		err = NotFound
	} else {
		run.PushScope(scope.AtFinder(data))
		ret, err = ps.GetNumber(run)
		run.PopScope()
	}
	return
}

func (m TextMap) GetTextMatching(run rt.Runtime, name string, data rt.Object) (ret string, err error) {
	id := id.MakeId(name)
	if ps, ok := m[id]; !ok {
		err = NotFound
	} else {
		run.PushScope(scope.AtFinder(data))
		ret, err = ps.GetText(run)
		run.PopScope()
	}
	return
}

func (m ObjectMap) GetObjectMatching(run rt.Runtime, name string, data rt.Object) (ret rt.Object, err error) {
	id := id.MakeId(name)
	if ps, ok := m[id]; !ok {
		err = NotFound
	} else {
		run.PushScope(scope.AtFinder(data))
		ret, err = ps.GetObject(run)
		run.PopScope()
	}
	return
}

func (m NumListMap) GetNumStreamMatching(run rt.Runtime, name string, data rt.Object) (ret rt.NumberStream, err error) {
	id := id.MakeId(name)
	if ps, ok := m[id]; !ok {
		err = NotFound
	} else {
		run.PushScope(scope.AtFinder(data))
		ret, err = ps.GetNumberStream(run)
		run.PopScope()
	}
	return
}

func (m TextListMap) GetTextStreamMatching(run rt.Runtime, name string, data rt.Object) (ret rt.TextStream, err error) {
	id := id.MakeId(name)
	if ps, ok := m[id]; !ok {
		err = NotFound
	} else {
		ret, err = ps.GetTextStream(run)
		run.PopScope()
	}
	return
}

func (m ObjListMap) GetObjStreamMatching(run rt.Runtime, name string, data rt.Object) (ret rt.ObjectStream, err error) {
	id := id.MakeId(name)
	if ps, ok := m[id]; !ok {
		err = NotFound
	} else {
		ret, err = ps.GetObjectStream(run)
		run.PopScope()
	}
	return
}
