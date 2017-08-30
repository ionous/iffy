package locate

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/index"
	"github.com/ionous/iffy/ref/class"
	"github.com/ionous/iffy/rt"
)

type Locale struct {
	*index.Table
}

func (l *Locale) Empty(p rt.Object) bool {
	_, hasChild := l.Primary.FindFirst(0, p.GetId().Name)
	return !hasChild

}
func (l *Locale) SetLocation(p, c rt.Object, now Containment) (err error) {
	types := map[Containment]struct {
		Parent, Child string
	}{
		Supports: {Parent: "container", Child: "thing"},
		Contains: {Parent: "container", Child: "thing"},
		Wears:    {Parent: "actor", Child: "thing"},
		Carries:  {Parent: "actor", Child: "thing"},
		Holds:    {Parent: "room", Child: "thing"},
	}
	if check, ok := types[now]; !ok {
		err = errutil.New("relation not supported", now)
	} else if !class.IsCompatible(p.GetClass(), check.Parent) {
		err = errutil.New("expected parent", check.Parent)
	} else if !class.IsCompatible(c.GetClass(), check.Child) {
		err = errutil.New("expected child", check.Child)
	} else {
		err = l.AddPair(p.GetId().Name, c.GetId().Name, func(old interface{}) (ret interface{}, err error) {
			if c, ok := old.(Containment); ok && c != now {
				err = errutil.New("was", c, "now", now)
			} else {
				ret = now
			}
			return
		})
	}
	return
}
