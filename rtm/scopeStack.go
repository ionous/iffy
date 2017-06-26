package rtm

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/reflector"
	"github.com/ionous/iffy/rt"
)

// ScopeStack implements rt.ObjectFinder by reflecting calls to the top of a stack of ObjectFinders.
type ScopeStack struct {
	scp []rt.ObjectFinder
}

// FindObject implements rt.ObjectFinder
func (os *ScopeStack) FindObject(name string) (ret ref.Object, okay bool) {
	if len(os.scp) > 0 {
		id := reflector.MakeId(name)
		ret, okay = os.scp[0].FindObject(id)
	}
	return
}

// PushScope activates the passed object finder.
func (os *ScopeStack) PushScope(scope rt.ObjectFinder) {
	os.scp = append(os.scp, scope)
}

// PopScope deactivates the object finder most recently passed to PushScope.
func (os *ScopeStack) PopScope() {
	if l := len(os.scp) - 1; l >= 0 {
		os.scp = os.scp[:l]
	} else {
		panic(errutil.New("cant pop empty scope"))
	}
}
