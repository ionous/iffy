package story

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/rt"
)

// run("activity", "{go+execute|ghost}");
func imp_activity(k *Importer, r reader.Map) (ret *core.Activity, err error) {
	var exes []rt.Execute
	if m, e := reader.Unpack(r, "activity"); e != nil {
		err = e
	} else if e := reader.Repeats(m.SliceOf("$EXE"),
		func(m reader.Map) (err error) {
			if i, e := k.DecodeSlot(m, "execute"); e != nil {
				err = e
			} else if i != nil {
				exes = append(exes, i.(rt.Execute))
			}
			return
		}); e != nil {
		err = e
	} else {
		ret = core.NewActivity(exes...)
	}
	return
}
