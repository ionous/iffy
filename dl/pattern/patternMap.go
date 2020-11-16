package pattern

import (
	r "reflect"

	"github.com/ionous/errutil"
)

// PatternMap: a simple helper, mainly for testing, to provide access to named patterns
type PatternMap map[string]interface{}

// skip assembling the pattern from the db
// we just want to test we can invoke a pattern successfully.
// pv is a pointer to a pattern instance, and we copy its contents in.
func (m PatternMap) GetEvalByName(name string, pv interface{}) (err error) {
	if patternPtr, ok := m[name]; !ok {
		err = errutil.New("unknown pattern", name)
	} else {
		stored := r.ValueOf(patternPtr).Elem()
		outVal := r.ValueOf(pv).Elem()
		outVal.Set(stored)
	}
	return
}
