package express

import (
	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/export"
)

func assignProps(out r.Value, args []r.Value) (err error) {
	export.WalkProperties(out.Type(), func(f *r.StructField, path []int) (done bool) {
		if len(args) <= 0 {
			done = true
		} else {
			var arg r.Value
			arg, args = args[0], args[1:]
			outAt := out.FieldByIndex(path)
			if argType := arg.Type(); !argType.AssignableTo(outAt.Type()) {
				err = errutil.New("cant assign %s to %q", argType.String(), f.Name)
			} else {
				outAt.Set(arg)
			}
		}
		return err != nil
	})
	if err == nil && len(args) > 0 {
		err = errutil.New("unable to consume all args")
	}
	return
}
