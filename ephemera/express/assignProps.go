package express

import (
	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/export"
)

func assignProps(out r.Value, args []r.Value) (err error) {
	outType := out.Type()
	export.WalkProperties(outType, func(f *r.StructField, path []int) (done bool) {
		if len(args) <= 0 {
			done = true
		} else {
			field := out.FieldByIndex(path)
			if f.Type.Kind() != r.Slice {
				arg := args[0]
				if argType := arg.Type(); !argType.AssignableTo(f.Type) {
					err = errutil.Fmt("cant assign %s to field %s{ %s %s }",
						argType, outType, f.Name, f.Type)
				} else {
					field.Set(arg)
					args = args[1:] // pop
				}
			} else {
				// when assigning to a slice, eat as many elements as possible.
				// it makes having slices as the last element of a command a good idea.
				slice, elType := field, f.Type.Elem()
				for len(args) > 0 {
					arg := args[0]
					if on, ok := arg.Interface().(objectName); ok {
						arg = r.ValueOf(on.getTextName())
					}
					if argType := arg.Type(); !argType.AssignableTo(elType) {
						break
					} else {
						slice = r.Append(slice, arg)
						args = args[1:] // pop
					}
				}
				field.Set(slice)
			}
		}
		return err != nil // returns "done" when there is an error.
	})
	if err == nil && len(args) > 0 {
		err = errutil.New("unable to consume all args")
	}
	return
}
