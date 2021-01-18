package decode

import (
	r "reflect"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/export"
	"github.com/ionous/iffy/export/tag"
)

// ReadRet is similar to reader.ReadMap, except it returns a value.
type ReadRet func(reader.Map) (interface{}, error)

type cmdRec struct {
	elem         r.Type
	customReader ReadRet
}

// Decoder reads programs from json.
type Decoder struct {
	source     string
	cmds       map[string]cmdRec
	issueFn    IssueReport
	IssueCount int
	Path       []string
}

func NewDecoder() *Decoder {
	reportNothing := func(reader.Position, error) {}
	return NewDecoderReporter("decoder", reportNothing)
}

// AddCallback registers a command parser.
func (dec *Decoder) AddCallback(cmd composer.Composer, cb ReadRet) {
	n := composer.SpecName(cmd)
	if was, exists := dec.cmds[n]; exists && was.customReader != nil {
		panic(errutil.Fmt("conflicting name for spec %q %q!=%T", n, was.elem, cmd))
	} else {
		elem := r.TypeOf(cmd).Elem()
		dec.cmds[n] = cmdRec{elem, cb}
	}
}

// AddDefaultCallbacks registers default command parsers.
func (dec *Decoder) AddDefaultCallbacks(slats []composer.Composer) {
	for _, cmd := range slats {
		dec.AddCallback(cmd, nil)
	}
}

func (dec *Decoder) ReadSpec(m reader.Map) (ret interface{}, err error) {
	if val, e := dec.importItem(m); e != nil {
		err = e
	} else {
		ret = val.Interface()
	}
	return
}

// returns a ptr r.Value
func (dec *Decoder) importItem(m reader.Map) (ret r.Value, err error) {
	itemValue, itemType := m, m.StrOf(reader.ItemType)
	if cmd, ok := dec.cmds[itemType]; !ok {
		err = errutil.Fmt("unknown type %q at %q", itemType, reader.At(m))
	} else {
		m := itemValue
		if custom := cmd.customReader; custom != nil {
			if res, e := custom(m); e != nil {
				err = e
			} else {
				ret = r.ValueOf(res)
			}
		} else if cmd.elem.Kind() != r.Struct {
			err = errutil.New("expected a struct", itemType, "is a", cmd.elem.String())
		} else {
			ptr := r.New(cmd.elem)
			dec.ReadFields(reader.At(m), ptr.Elem(), m.MapOf(reader.ItemValue))
			ret = ptr
		}
	}
	return
}

var posType = r.TypeOf((*reader.Position)(nil)).Elem()

func (dec *Decoder) ReadFields(at string, out r.Value, in reader.Map) {
	name := out.Type().String()
	dec.Path = append(dec.Path, name)
	//
	var fields []string
	export.WalkProperties(out.Type(), func(f *r.StructField, path []int) (done bool) {
		token := export.Tokenize(f.Name)
		fields = append(fields, token)
		// we report on missing properties below.
		if inVal, ok := in[token]; !ok {
			// log only if the field is required. not optional.
			if t := tag.ReadTag(f.Tag); t.Exists("internal") {
				if f.Type == posType {
					outAt := out.FieldByIndex(path)
					outAt.Set(r.ValueOf(reader.Position{dec.source, at}))
				}
			} else if !t.Exists("optional") {
				// and even then only if its a fixed field
				if f.Type.Kind() != r.Ptr {
					dec.report(at, errutil.Fmt("missing %s.%s at %s", out.Type().String(), token, at))
				}
			}
		} else {
			outAt := out.FieldByIndex(path)
			if e := dec.importValue(outAt, inVal); e != nil {
				dec.report(at, errutil.New("error processing field", out.Type().String(), f.Name, e))
			}
		}
		return
	})
	dec.Path = dec.Path[0 : len(dec.Path)-1]

	// walk keys of json dictionary:
	for token, _ := range in {
		i, cnt := 0, len(fields)
		for ; i < cnt; i++ {
			if token == fields[i] {
				break
			}
		}
		if i == cnt {
			dec.report(at, errutil.Fmt("unprocessed %q", token))
		}
	}
}

func (dec *Decoder) importValue(outAt r.Value, inVal interface{}) (err error) {
	switch outType := outAt.Type(); outType.Kind() {
	case r.Float32, r.Float64:
		err = dec.unpack(inVal, func(p reader.Map, v interface{}) (err error) {
			// float64, for JSON numbers
			if n, ok := v.(float64); !ok {
				err = errutil.Fmt("expected a number, have %T", v)
			} else {
				outAt.SetFloat(n)
			}
			return
		})
	case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
		err = dec.unpack(inVal, func(p reader.Map, v interface{}) (err error) {
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
		err = dec.unpack(inVal, func(p reader.Map, v interface{}) (err error) {
			// string, for JSON strings
			if str, ok := v.(string); !ok {
				err = errutil.New("expected a string")
			} else {
				outAt.SetBool(str == "$TRUE") // only need to set true: false is the zero value.
			}
			return
		})

	case r.String:
		err = dec.unpack(inVal, func(p reader.Map, v interface{}) (err error) {
			// string, for JSON strings
			if str, ok := v.(string); !ok {
				err = errutil.New("expected a string")
			} else {
				outAt.SetString(str)
			}
			return
		})

	case r.Ptr:
		// see if its an optional value.
		ptr := r.New(outAt.Type().Elem())
		if e := dec.importValue(ptr.Elem(), inVal); e != nil {
			err = e
		} else {
			outAt.Set(ptr)
		}

	case r.Struct:
		// b/c of the way optional values are specified,
		// going from r.Struct is easier than from r.Ptr.
		switch spec := outAt.Addr().Interface().(type) {
		case StrType:
			if e := dec.unpack(inVal, func(p reader.Map, v interface{}) (err error) {
				if str, ok := v.(string); !ok {
					err = errutil.Fmt("expected string, got %T(%v)", v, v)
				} else {
					// fix?: by using field by name we "unwrap" embedded structs
					// ex. VariableName { core.Variable }
					outAt.FieldByName("Str").SetString(str)
					storeAt(p, outAt)
				}
				return
			}); e != nil {
				err = e
			}

		case NumType:
			if e := dec.unpack(inVal, func(p reader.Map, v interface{}) (err error) {
				if num, ok := v.(float64); !ok {
					err = errutil.Fmt("expected float, got %T(%v)", v, v)
				} else {
					// validate choice; fix: tolerance?
					// handle conversion b/t floats and ints of different widths
					tgt := outAt.Field(outAt.NumField() - 1)
					v := r.ValueOf(num).Convert(tgt.Type())
					tgt.Set(v)
				}
				return
			}); e != nil {
				err = e
			}

		case SwapType:
			if e := dec.unpack(inVal, func(p reader.Map, v interface{}) (err error) {
				if data, ok := v.(map[string]interface{}); !ok {
					err = errutil.Fmt("expected swap, got %T(%v)", v, v)
				} else {
					found := false
					for k, typePtr := range spec.Choices() {
						token := "$" + strings.ToUpper(k)
						if contents, ok := data[token]; ok {
							ptr := r.New(r.TypeOf(typePtr).Elem())
							if e := dec.importValue(ptr.Elem(), contents); e != nil {
								err = e
							} else {
								storeAt(p, outAt)
								outAt.Field(outAt.NumField() - 1).Set(ptr)
							}
							found = true
							break
						}
					}
					if !found {
						err = errutil.New("no valid swap data found")
					}
				}
				return
			}); e != nil {
				err = e
			}

		default:
			if e := dec.unpack(inVal, func(p reader.Map, v interface{}) (err error) {
				if slot, ok := v.(map[string]interface{}); !ok {
					err = errutil.Fmt("expected map for %T, got %T(%v)", spec, v, v)
				} else {
					dec.ReadFields(reader.At(p), outAt, reader.Map(slot))
				}
				return
			}); e != nil {
				err = e
			}
		}

	case r.Interface:
		if e := dec.unpack(inVal, func(p reader.Map, v interface{}) (err error) {
			// map[string]interface{}, for JSON objects
			if slot, ok := v.(map[string]interface{}); !ok {
				err = errutil.Fmt("expected map for interface, got %T(%v)", v, v)
			} else if val, e := dec.importItem(slot); e != nil {
				dec.report(reader.At(slot), e)
			} else if rtype := val.Type(); !rtype.AssignableTo(outAt.Type()) {
				err = errutil.New(rtype, "not assignable to", outType)
			} else {
				outAt.Set(val)
			}
			return
		}); e != nil {
			err = e
		}

	case r.Slice:
		// []interface{}, for JSON arrays
		if items, ok := inVal.([]interface{}); !ok {
			err = errutil.New("expected a slice")
		} else {
			elType := outType.Elem()
			if slice := outAt; len(items) > 0 {
				for _, item := range items {
					if k := elType.Kind(); k != r.Interface {
						el := r.New(elType).Elem()
						if e := dec.importValue(el, item); e != nil {
							err = errutil.Append(err, e)
						} else {
							slice = r.Append(slice, el)
						}
					} else {
						// note: this skips over the slot itself ( ex execute )
						if e := dec.unpack(item, func(p reader.Map, v interface{}) (err error) {
							// map[string]interface{}, for JSON objects
							if itemData, ok := v.(map[string]interface{}); !ok {
								// execute has some single nulls sometimes;
								if v != nil {
									err = errutil.Fmt("item data not a slot %T", itemData)
								}
							} else if val, e := dec.importItem(itemData); e != nil {
								err = e // elType is ex. *story.Paragraph; itemData has a member $STORY_STATEMENT
							} else if rtype := val.Type(); !rtype.AssignableTo(elType) {
								err = errutil.New(rtype, "not assignable to", elType)
							} else {
								slice = r.Append(slice, val)
							}
							return
						}); e != nil {
							err = errutil.Append(err, e)
						}
					}
				}
				outAt.Set(slice)
			}
		}
	}
	return
}

func storeAt(m reader.Map, val r.Value) {
	if at := reader.At(m); len(at) > 0 {
		if v := val.FieldByName("At"); !v.IsValid() {
			v.Set(r.ValueOf(at))
		}
	}
}

// cast inVal to a map, and call setter with contents of "value"
func (dec *Decoder) unpack(inVal interface{}, setter func(reader.Map, interface{}) error) (err error) {
	if item, ok := inVal.(map[string]interface{}); !ok {
		err = errutil.New("expected an item, got:", inVal)
	} else {
		val := item[reader.ItemValue]
		if e := setter(item, val); e != nil {
			dec.report(reader.At(item), e)
		}
	}
	return
}
