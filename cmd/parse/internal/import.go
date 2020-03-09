package internal

import (
	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/export"
	"github.com/ionous/iffy/ref/kindOf"
	"github.com/ionous/iffy/ref/unique"
)

// -----------------------------------
const (
	itemId    = "id"
	itemType  = "type"
	itemValue = "value"
)

func Import(targetPtr interface{}, inData export.Dict, types map[string]export.Run) (err error) {
	out := r.ValueOf(targetPtr).Elem()
	return Unmarshall(out, inData, types)
}

func Unmarshall(out r.Value, inData export.Dict, types map[string]export.Run) (err error) {
	if inVal, ok := inData[itemValue].(map[string]interface{}); !ok {
		err = errutil.New("unexpected value in data", inData)
	} else if e := unmarshall(out, inVal, types); e != nil {
		id, _ := inData[itemId].(string)
		err = errutil.Append(errutil.New("Unmarshall", id, "error(s):"), e)
	}
	return
}

func unmarshall(out r.Value, in export.Dict, types map[string]export.Run) (err error) {
	var processed []string

	unique.WalkProperties(out.Type(), func(f *r.StructField, path []int) (done bool) {
		token := export.Tokenize(f)
		processed = append(processed, token)

		// value for the field not found? that's okay.
		// note: values of run-fields are always going to be an "item" or an array of items
		if inVal, ok := in[token]; ok {
			outAt := out.FieldByIndex(path)
			if e := importValue(outAt, inVal, types); e != nil {
				e := errutil.New("error processing field", f.Name, e)
				err = errutil.Append(err, e)
			}
		}
		return
	})

	// walk keys of json dictionary:
	for token, _ := range in {
		found := false
		for _, key := range processed {
			if key == token {
				found = true
				break
			}
		}
		if !found {
			e := errutil.New("unprocessed value", token)
			err = errutil.Append(err, e)
		}
	}
	return
}

func importValue(outAt r.Value, inVal interface{}, types map[string]export.Run) (err error) {
	switch outType := outAt.Type(); {
	case kindOf.Float(outType):
		err = unpack(inVal, func(v interface{}) (err error) {
			// float64, for JSON numbers
			if n, ok := v.(float64); !ok {
				err = errutil.New("expected a number")
			} else {
				outAt.SetFloat(n)

			}
			return
		})
	case kindOf.Int(outType):
		err = unpack(inVal, func(v interface{}) (err error) {
			// float64, for JSON numbers
			if n, ok := v.(float64); !ok {
				err = errutil.New("expected a number")
			} else {
				outAt.SetInt(int64(n))
			}
			return
		})

	case outType.Kind() == r.Bool:
		// fix? boolean values are stored as enumerations
		err = unpack(inVal, func(v interface{}) (err error) {
			// string, for JSON strings
			if str, ok := v.(string); !ok {
				err = errutil.New("expected a number")
			} else {
				outAt.SetBool(str == "$TRUE") // only need to set true: false is the zero value.
			}
			return
		})

	case outType.Kind() == r.String:
		err = unpack(inVal, func(v interface{}) (err error) {
			// string, for JSON strings
			if str, ok := v.(string); !ok {
				err = errutil.New("expected a number")
			} else {
				outAt.SetString(str)
			}
			return
		})

	case outType.Kind() == r.Interface:
		if e := unpack(inVal, func(v interface{}) (err error) {
			// map[string]interface{}, for JSON objects
			if slot, ok := v.(map[string]interface{}); !ok {
				err = errutil.New("value not a slot")
			} else if v, e := importSlot(slot, outAt.Type(), types); e != nil {
				err = e
			} else {
				outAt.Set(v)
			}
			return
		}); e != nil {
			err = errutil.Append(err, e)
		}

	case outType.Kind() == r.Slice && kindOf.Interface(outType.Elem()):
		// []interface{}, for JSON arrays
		if items, ok := inVal.([]interface{}); ok {
			elType := outType.Elem()
			if slice := outAt; len(items) > 0 {
				for _, item := range items {
					if e := unpack(item, func(v interface{}) (err error) {
						// map[string]interface{}, for JSON objects
						if slot, ok := v.(map[string]interface{}); !ok {
							err = errutil.New("value not a slot")
						} else if v, e := importSlot(slot, elType, types); e != nil {
							err = e
						} else {
							slice = r.Append(slice, v)
						}
						return
					}); e != nil {
						err = errutil.Append(err, e)
					}
				}
				outAt.Set(slice)
			}
		}
	}
	return
}

func unpack(inVal interface{}, setter func(interface{}) error) (err error) {
	if item, ok := inVal.(map[string]interface{}); !ok {
		err = errutil.New("expected an item, got:", inVal)
	} else if e := setter(item[itemValue]); e != nil {
		id, _ := item[itemId].(string)
		err = errutil.New("couldnt unpack", id, e)
	}
	return
}

func importSlot(slot export.Dict, slotType r.Type, types map[string]export.Run) (ret r.Value, err error) {
	typeName, _ := slot[itemType].(string)
	if run, ok := types[typeName]; !ok {
		err = errutil.New("unknown type", typeName, slot)
	} else {
		rtype := r.TypeOf(run.Type)
		if !rtype.AssignableTo(slotType) {
			err = errutil.New("incompatible types", rtype.String(), "not assignable to", slotType.String())
		} else {
			v := r.New(rtype.Elem())
			if e := Unmarshall(v.Elem(), slot, types); e != nil {
				err = e
			} else {
				ret = v
			}
		}
	}
	return
}
