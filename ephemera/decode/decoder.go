package decode

import (
	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/export"
)

// ReadRet is similar to reader.ReadMap, except it returns a value.
type ReadRet func(reader.Map) (interface{}, error)

type Override struct {
	Spec     composer.Slat
	Callback ReadRet
}

// Decoder reads programs from json.
type Decoder struct {
	cmds map[string]ReadRet
}

func NewDecoder() *Decoder {
	dec := &Decoder{make(map[string]ReadRet)}
	return dec
}

func (dec *Decoder) AddCallbacks(overrides []Override) {
	for _, n := range overrides {
		dec.AddCallback(n.Spec, n.Callback)
	}
}

// AddCallback registers a command parser.
func (dec *Decoder) AddCallback(cmd composer.Slat, cb ReadRet) {
	if spec := cmd.Compose(); len(spec.Name) == 0 {
		panic(errutil.Fmt("missing name for spec %T", cmd))
	} else {
		dec.cmds[spec.Name] = cb
	}
}

// AddDefaultCallbacks registers default command parsers.
func (dec *Decoder) AddDefaultCallbacks(slats []composer.Slat) {
	for _, slat := range slats {
		spec := slat.Compose()
		elem := r.TypeOf(slat).Elem()
		dec.cmds[spec.Name] = func(m reader.Map) (ret interface{}, err error) {
			return dec.readNew(m, elem)
		}
	}
}

// ReadProg attempts to parse the passed json data as a golang program.
func (dec *Decoder) ReadProg(m reader.Map) (ret interface{}, err error) {
	itemValue, itemType := m, m.StrOf(reader.ItemType)
	if fn, ok := dec.cmds[itemType]; !ok {
		err = errutil.Fmt("unknown type %q", itemType)
	} else {
		ret, err = fn(itemValue)
	}
	return
}

// m is the contents of slotType is a concrete command ( not a ptr to a command )
func (dec *Decoder) readNew(m reader.Map, slotType r.Type) (ret interface{}, err error) {
	if slotType.Kind() != r.Struct {
		panic("expected a struct")
	} else {
		ptr := r.New(slotType)
		if e := dec.readFields(ptr.Elem(), m.MapOf(reader.ItemValue)); e != nil {
			err = e
		} else {
			ret = ptr.Interface()
		}
	}
	return
}

func (dec *Decoder) readFields(out r.Value, in reader.Map) (err error) {
	var processed []string
	export.WalkProperties(out.Type(), func(f *r.StructField, path []int) (done bool) {
		token := export.Tokenize(f)
		processed = append(processed, token)
		// value for the field not found? that's okay.
		// note: values of run-fields are always going to be an "item" or an array of items
		if inVal, ok := in[token]; ok {
			outAt := out.FieldByIndex(path)
			if e := dec.importValue(outAt, inVal); e != nil {
				e := errutil.New("error processing field", out.Type().String(), f.Name, e)
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
			e := errutil.Fmt("unprocessed value %q", token)
			err = errutil.Append(err, e)
		}
	}
	return
}

// returns a ptr r.Value
func (dec *Decoder) importSlot(m reader.Map, slotType r.Type) (ret r.Value, err error) {
	itemValue, itemType := m, m.StrOf(reader.ItemType)
	if cmdImport, ok := dec.cmds[itemType]; !ok {
		err = errutil.New("unknown type", itemType, reader.At(m))
	} else if cmd, e := cmdImport(itemValue); e != nil {
		err = e
	} else {
		rval := r.ValueOf(cmd)
		if rtype := rval.Type(); !rtype.AssignableTo(slotType) {
			err = errutil.New("incompatible types", rtype.String(), "not assignable to", slotType.String())
		} else {
			ret = rval
		}
	}
	return
}

func (dec *Decoder) importValue(outAt r.Value, inVal interface{}) (err error) {
	switch outType := outAt.Type(); outType.Kind() {
	case r.Float32, r.Float64:
		err = unpack(inVal, func(v interface{}) (err error) {
			// float64, for JSON numbers
			if n, ok := v.(float64); !ok {
				err = errutil.Fmt("expected a number, have %T", v)
			} else {
				outAt.SetFloat(n)
			}
			return
		})
	case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
		err = unpack(inVal, func(v interface{}) (err error) {
			// float64, for JSON numbers
			if n, ok := v.(float64); !ok {
				err = errutil.New("expected a number")
			} else {
				outAt.SetInt(int64(n))
			}
			return
		})

	case r.Bool:
		// fix? boolean values are stored as enumerations
		err = unpack(inVal, func(v interface{}) (err error) {
			// string, for JSON strings
			if str, ok := v.(string); !ok {
				err = errutil.New("expected a string")
			} else {
				outAt.SetBool(str == "$TRUE") // only need to set true: false is the zero value.
			}
			return
		})

	case r.String:
		err = unpack(inVal, func(v interface{}) (err error) {
			// string, for JSON strings
			if str, ok := v.(string); !ok {
				err = errutil.New("expected a string")
			} else {
				outAt.SetString(str)
			}
			return
		})

	case r.Interface:
		// note: this skips over the slot itself ( ex execute )
		if e := unpack(inVal, func(v interface{}) (err error) {
			// map[string]interface{}, for JSON objects
			if slot, ok := v.(map[string]interface{}); !ok {
				err = errutil.New("value not a slot")
			} else if v, e := dec.importSlot(slot, outAt.Type()); e != nil {
				err = e
			} else {
				outAt.Set(v)
			}
			return
		}); e != nil {
			err = errutil.Append(err, e)
		}

	case r.Slice:
		if outType.Elem().Kind() == r.Interface {
			// []interface{}, for JSON arrays
			if items, ok := inVal.([]interface{}); ok {
				elType := outType.Elem()
				if slice := outAt; len(items) > 0 {
					for _, item := range items {
						// note: this skips over the slot itself ( ex execute )
						if e := unpack(item, func(v interface{}) (err error) {
							// map[string]interface{}, for JSON objects
							if itemData, ok := v.(map[string]interface{}); !ok {
								err = errutil.Fmt("item data not a slot %T", itemData)
							} else if v, e := dec.importSlot(itemData, elType); e != nil {
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
	}
	return
}

func unpack(inVal interface{}, setter func(interface{}) error) (err error) {
	if item, ok := inVal.(map[string]interface{}); !ok {
		err = errutil.New("expected an item, got:", inVal)
	} else if e := setter(item[reader.ItemValue]); e != nil {
		id, _ := item[reader.ItemId].(string)
		err = errutil.New("couldnt unpack", id, e)
	}
	return
}
