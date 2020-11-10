package generic

// CopyFloats: duplicate the passed slice.
// ( b/c golang's built in copy doesnt allocate )
func CopyFloats(src []float64) []float64 {
	out := make([]float64, len(src))
	copy(out, src)
	return out
}

// CopyStrings: duplicate the passed slice.
// ( b/c golang's built in copy doesnt allocate )
func CopyStrings(src []string) []string {
	out := make([]string, len(src))
	copy(out, src)
	return out
}

// CopyValue: create a new generic value capable of supporting the passed affinity.
// from a snapshot of the passed value; errors if the two types are not compatible.
// func CopyValue(kinds Kinds, val Value) (ret Value, err error) {
// 	if val == nil {
// 		err = errutil.New("failed to copy nil value")
// 	} else {
// 		switch a := val.Affinity(); a {
// 		case affine.Bool:
// 			if v, e := val.GetBool(); e != nil {
// 				err = e
// 			} else {
// 				ret = BoolOf(v)
// 			}
// 		case affine.Number:
// 			if v, e := val.GetNumber(); e != nil {
// 				err = e
// 			} else {
// 				ret = FloatOf(v)
// 			}
// 		case affine.NumList:
// 			if vs, e := val.GetNumList(); e != nil {
// 				err = e
// 			} else {
// 				ret = FloatsOf(CopyFloats(vs))
// 			}

// 		case affine.Text:
// 			if v, e := val.GetText(); e != nil {
// 				err = e
// 			} else {
// 				ret = StringOf(v)
// 			}
// 		case affine.TextList:
// 			if vs, e := val.GetTextList(); e != nil {
// 				err = e
// 			} else {
// 				ret = StringsOf(CopyStrings(vs))
// 			}

// 		case affine.Record:
// 			// could also peek under the hood by casting to .(*Record)
// 			if kind, e := kinds.GetKindByName(val.Type()); e != nil {
// 				err = errutil.New("unknown kind", val.Type(), e)
// 			} else if next, e := copyRecord(kind, val); e != nil {
// 				err = e
// 			} else {
// 				ret = next
// 			}
// 		case affine.RecordList:
// 			if kind, e := kinds.GetKindByName(val.Type()); e != nil {
// 				err = errutil.New("unknown kind", val.Type(), e)
// 			} else if cnt, e := val.GetLen(); e != nil {
// 				err = e
// 			} else {
// 				values := make([]*Record, cnt)
// 				for i := 0; i < cnt; i++ {
// 					if el, e := val.GetIndex(i); e != nil {
// 						err = e
// 						break
// 					} else if cpy, e := copyRecord(kind, el); e != nil {
// 						err = e
// 						break
// 					} else {
// 						values[i] = cpy
// 					}
// 				}
// 				if err == nil {
// 					ret = &RecordSlice{kind: kind, values: values}
// 				}
// 			}
// 		case affine.Object:
// 			// new nouns cant be dynamically added to the runtime.
// 			err = errutil.New("can't duplicate object values")

// 		default:
// 			err = errutil.Fmt("failed to copy value, expected %s got %v(%T)", a, val, val)
// 		}
// 	}
// 	return
// }

// // assumes in value is a record.
// func copyRecord(kind *Kind, val Value) (ret *Record, err error) {
// 	cnt := kind.NumField()
// 	values := make([]Value, cnt)
// 	for i := 0; i < cnt; i++ {
// 		ft := kind.Field(i) // fix: get field by index?
// 		if el, e := val.GetNamedField(ft.Name); e != nil {
// 			err = e
// 			break
// 			if cpy, e := CopyValue(kind.kinds, el); e != nil {
// 				err = e
// 				break
// 			} else {
// 				values[i] = cpy
// 			}
// 		}
// 	}
// 	if err == nil {
// 		ret = &Record{kind: kind, values: values}
// 	}
// 	return
// }
