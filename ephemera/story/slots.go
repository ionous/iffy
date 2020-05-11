package story

import (
	"reflect"
	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/export"
)

// FIX: some sort of reverse lookup instead?
func slotName(i interface{}) (ret string, err error) {
	itype := r.TypeOf(i)
	found := false
	for _, slot := range export.Slots {
		rtype := reflect.TypeOf(slot.Type).Elem()
		if itype.Implements(rtype) {
			ret = slot.Name
			found = true
			break
		}
	}
	if !found {
		err = errutil.New("couldnt determine matching slot %T", i)
	}
	return
}
