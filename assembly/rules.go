package assembly

import (
	"bytes"
	"encoding/gob"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/tables"
)

// the first parameter should be a *string, the second some *bytes
func (b *BuildRule) buildFromRule(asm *Assembler, args ...interface{}) (err error) {
	list := make(map[string]interface{})
	var last string
	var curr interface{}
	if e := tables.QueryAll(asm.cache.DB(), b.Query,
		func() (err error) {
			name, prog := *args[0].(*string), *args[1].(*[]byte)
			if name != last || curr == nil {
				curr = b.NewContainer(name)
				list[name] = curr
				last = name
			}
			el := b.NewEl(curr)
			dec := gob.NewDecoder(bytes.NewBuffer(prog))
			return dec.Decode(el)
		}, args...); e != nil {
		err = errutil.New("buildFromRule", e)
	} else {
		// write the passed list of gobs into the assembler db
		err = asm.WriteGobs(list)
	}
	return
}
