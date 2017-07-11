package rtm

import (
	"github.com/ionous/iffy/pat"
	"github.com/ionous/iffy/rt"
)

// Thunk from the rt.Patterns interface to the pat.Patterns interface
// FIX? remove "scope" setup from package pat, and move to here?
type Thunk struct {
	run rt.Runtime
	pat.Patterns
}

func (m Thunk) GetBoolMatching(data rt.Object) (bool, error) {
	return m.Patterns.GetBoolMatching(m.run, data)
}

func (m Thunk) GetNumMatching(data rt.Object) (float64, error) {
	return m.Patterns.GetNumMatching(m.run, data)
}

func (m Thunk) GetTextMatching(data rt.Object) (string, error) {
	return m.Patterns.GetTextMatching(m.run, data)
}

func (m Thunk) GetObjectMatching(data rt.Object) (rt.Object, error) {
	return m.Patterns.GetObjectMatching(m.run, data)
}

func (m Thunk) GetNumStreamMatching(data rt.Object) (rt.NumberStream, error) {
	return m.Patterns.GetNumStreamMatching(m.run, data)
}

func (m Thunk) GetTextStreamMatching(data rt.Object) (rt.TextStream, error) {
	return m.Patterns.GetTextStreamMatching(m.run, data)
}

func (m Thunk) GetObjStreamMatching(data rt.Object) (rt.ObjectStream, error) {
	return m.Patterns.GetObjStreamMatching(m.run, data)
}

func (m Thunk) ExecuteMatching(data rt.Object) (bool, error) {
	return m.Patterns.ExecuteMatching(m.run, data)
}
