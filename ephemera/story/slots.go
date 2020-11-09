package story

import (
	"reflect"
	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy"
)

// FIX: some sort of reverse lookup instead?
// ex. see findTypeName... maybe share via exports?
func slotName(i interface{}) (ret string, err error) {
	if i != nil { // null happens for FromVar
		itype := r.TypeOf(i)
		found := false
		for _, slots := range iffy.AllSlots {
			for _, slot := range slots {
				rtype := reflect.TypeOf(slot.Type).Elem()
				if itype.Implements(rtype) {
					ret = slot.Name
					found = true
					break
				}
			}
		}
		if !found {
			err = errutil.Fmt("couldnt determine matching slot %T", i)
		}
	}
	return
}
