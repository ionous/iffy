package composer

import (
	r "reflect"

	"github.com/ionous/iffy/lang"
)

// fix? imposes some otherwise unneeded imports....
func SpecName(c Composer) (ret string) {
	if c == nil {
		ret = "<nil>"
	} else if spec := c.Compose(); len(spec.Name) > 0 {
		ret = spec.Name
	} else {
		el := r.TypeOf(c).Elem()
		ret = lang.Underscore(el.Name())
	}
	return
}

func SlotName(c Slot) (ret string) {
	if n := c.Name; len(n) > 0 {
		ret = n
	} else {
		el := r.TypeOf(c).Elem()
		ret = lang.Underscore(el.Name())
	}
	return
}
