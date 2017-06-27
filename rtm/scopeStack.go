package rtm

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

// ScopeStack implements rt.ObjectFinder by reflecting calls to the top of a stack of ObjectFinders.
type ScopeStack struct {
	stack []rt.ObjectFinder
}

// FindObject implements rt.ObjectFinder
func (os *ScopeStack) FindObject(name string) (ret rt.Object, okay bool) {
	if cnt := len(os.stack); cnt > 0 {
		ret, okay = os.stack[cnt-1].FindObject(name)
	}
	return
}

// PushScope activates the passed object finder.
func (os *ScopeStack) PushScope(scope rt.ObjectFinder) {
	os.stack = append(os.stack, scope)
}

// PopScope deactivates the object finder most recently passed to PushScope.
func (os *ScopeStack) PopScope() {
	if cnt := len(os.stack); cnt > 0 {
		os.stack = os.stack[:cnt-1]
	} else {
		panic(errutil.New("cant pop empty scope"))
	}
}
