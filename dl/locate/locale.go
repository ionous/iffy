package locate

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/index"
	"github.com/ionous/iffy/ref/class"
	"github.com/ionous/iffy/rt"
)

type Locale struct {
	*index.Table // one to many
}

func (l *Locale) Empty(p rt.Object) bool {
	_, hasChild := l.Primary.FindFirst(0, p.Id().Name)
	return !hasChild
}

func (l *Locale) SetLocation(p rt.Object, now Containment, c rt.Object) (err error) {
	if check, ok := types[now]; !ok {
		err = errutil.New("relation not supported", now)
	} else if !class.IsCompatible(p.Type(), check.Parent) {
		err = errutil.New("expected parent", check.Parent)
	} else if !class.IsCompatible(c.Type(), check.Child) {
		err = errutil.New("expected child", check.Child)
	} else {
		err = l.AddPair(p.Id().Name, c.Id().Name, func(old interface{}) (ret interface{}, err error) {
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

var types = map[Containment]struct {
	Parent, Child string
}{
	Supports: {Parent: "container", Child: "thing"},
	Contains: {Parent: "container", Child: "thing"},
	Wears:    {Parent: "actor", Child: "thing"},
	Carries:  {Parent: "actor", Child: "thing"},
	Has:      {Parent: "room", Child: "thing"},
}
