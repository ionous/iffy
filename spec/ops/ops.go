package ops

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core" // for literals/literally.
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec"
	"github.com/ionous/iffy/spec/builder"
	r "reflect"
)

type Ops struct {
	unique.Types
	ShadowTypes unique.Types
}

func NewOps() *Ops {
	return &Ops{make(unique.Types), make(unique.Types)}
}

type Builder struct {
	builder.Builder
}

// NewBuilder starts creating a call tree.
func (ops *Ops) NewBuilder(root interface{}) (*Builder, bool) {
	target := r.ValueOf(root).Elem()
	spec := &_Spec{cmds: ops, target: target}
	return &Builder{
		builder.NewBuilder(&_Factory{ops}, spec),
	}, true
}

func (u *Builder) Build() (err error) {
	if _, e := u.Builder.Build(); e != nil {
		err = e
	}
	return
}

type _Factory struct {
	cmds *Ops
}

// _Target handles the differences between command structs and constructors.
// Commands are implemented by reflect.Value, and _Target is compatible with reflect's interface.
// FIX? support aggregate ops by using unique.Fields(), FieldByIndex()
type _Target interface {
	Type() r.Type
	Field(int) r.Value // Field returns the value of the requsted field. To maintain compatibility with reflect.Value: on error, Field returns an invalid Value.
	Addr() r.Value
}

// NewSpec implements spec.SpecFactory.
func (fac *_Factory) NewSpec(name string) (ret spec.Spec, err error) {
	if rtype, ok := fac.cmds.FindType(name); ok {
		target := r.New(rtype).Elem()
		ret = &_Spec{
			cmds:   fac.cmds,
			target: target,
		}
	} else if rtype, ok := fac.cmds.ShadowTypes.FindType(name); ok {
		shadow := &ShadowClass{rtype, make(map[string]_ShadowSlot)}
		ret = &_Spec{
			cmds:   fac.cmds,
			target: shadow,
		}
	} else {
		err = errutil.New("unknown command", name)
	}
	return
}

// NewSpecs implements spec.SpecFactory.
// the spec algorithm creates NewSpecs, and then assigns it to a slot
// we need the slot to target the array properly, so we just wait,
func (fac *_Factory) NewSpecs() (spec.Specs, error) {
	return &_Specs{cmds: fac.cmds}, nil
}

type _Spec struct {
	cmds   *Ops
	target _Target // output object we are building
	index  int
}

func (spec *_Spec) Position(arg interface{}) (err error) {
	tgt := spec.target
	tgtType := tgt.Type()
	if cnt := tgtType.NumField(); spec.index >= cnt {
		err = errutil.New("too many arguments", tgtType, "expected", cnt)
	} else {
		field := tgt.Field(spec.index)
		if !field.IsValid() {
			fieldName := tgtType.Field(spec.index).Name
			err = errutil.Fmt("couldnt get value for position %d (%s.%s) %v", spec.index, tgtType, fieldName, arg)
		} else if e := setField(field, arg); e != nil {
			fieldName := tgtType.Field(spec.index).Name
			err = errutil.Fmt("position %d (%s.%s) %v", spec.index, tgtType, fieldName, e)
		} else {
			spec.index++
		}
	}
	return
}

func (spec *_Spec) Assign(key string, arg interface{}) (err error) {
	myid := id.MakeId(key)
	tgt := spec.target
	tgtType := tgt.Type()
	// FIX: maybe via unique.Fields() -- almost like a RefClass? -- build a cache of id->field index. searching every assign is annoying.
	for i, cnt := spec.index, tgtType.NumField(); i < cnt; i++ {
		fieldInfo := tgtType.Field(i)
		if myid == id.MakeId(fieldInfo.Name) {
			field := tgt.Field(i)
			if !field.IsValid() {
				err = errutil.New("couldnt get value for", tgtType, fieldInfo.Name)
			} else if e := setField(field, arg); e != nil {
				err = errutil.New("field", key, e)
			}
			break
		}
	}
	return
}

type _Specs struct {
	cmds *Ops
	els  []*_Spec
}

func (specs *_Specs) AddElement(el spec.Spec) (err error) {
	if spec, ok := el.(*_Spec); !ok {
		err = errutil.Fmt("unexpected element type %T", el)
	} else {
		specs.els = append(specs.els, spec)
	}
	return
}

// dst is the field we are setting
func setField(dst r.Value, src interface{}) (err error) {
	switch src := src.(type) {
	case *_Spec:
		// all commands are interfaces are implemented with pointers
		targetPtr := src.target.Addr()
		if e := ref.CoerceValue(dst, targetPtr); e != nil {
			err = errutil.New("couldnt assign command", e)
		}

	case *_Specs:
		if kind, isArray := arrayKind(dst.Type()); !isArray || kind != r.Interface {
			if !isArray {
				err = errutil.Fmt("trying to set an array to %v", dst.Type())
			} else {
				err = errutil.New("trying to set commands to", kind)
			}
		} else {
			slice, elType := dst, dst.Type().Elem()
			for _, spec := range src.els {
				// all commands are interfaces are implemented with pointers
				rvalue := spec.target.Addr()
				if from := rvalue.Type(); !from.AssignableTo(elType) {
					err = errutil.Fmt("incompatible element type. from: %v to: %v", from, elType)
					break
				} else {
					slice = r.Append(slice, rvalue)
				}
			}
			dst.Set(slice)
		}
	default:
		if dst.Kind() == r.Interface {
			if literal, ok := literally(dst.Type(), src); ok {
				src = literal
			}
		}
		if e := ref.CoerceValue(dst, src); e != nil {
			err = errutil.New("couldnt assign value", e)
		}
	}
	return
}

// literally allows users to specify primitive values for some evals.
//
// c.Cmd("texts", sliceOf.String("one", "two", "three"))
// c.Value(sliceOf.String("one", "two", "three"))
//
// c.Cmd("get").Begin() { c.Cmd("object", "@") c.Value("text") }
// c.Cmd("get", "@", "text")
//
// FIX? move literals to "builtin" to avoid the dependency on core.
// ( or, more ugly, have a "shortcut" interface users of core can inject. )
func literally(dstType r.Type, src interface{}) (ret interface{}, okay bool) {
	switch src := src.(type) {
	case bool:
		ret, okay = &core.Bool{src}, true
	case float64:
		ret, okay = &core.Num{src}, true
	case []float64:
		ret, okay = &core.Numbers{src}, true
	case string:
		// could be text or object --
		switch dstType {
		case textEval:
			ret, okay = &core.Text{src}, true
		case objEval:
			ret, okay = &core.Object{src}, true
		}
	case []string:
		switch dstType {
		case textListEval:
			ret, okay = &core.Texts{src}, true
		case objListEval:
			ret, okay = &core.Objects{src}, true
		}
	default:
		// FIX? a cleaner way to convert any number like thing to float64?
		if v := r.ValueOf(src); numberKind(v.Kind()) {
			var num float64
			if e := ref.CoerceValue(r.ValueOf(&num).Elem(), v); e == nil {
				ret, okay = &core.Num{num}, true
			}
		}
	}
	return
}

// determine what kind of eval can produce the passed type.
func evalFromType(rtype r.Type) (ret r.Type, okay bool) {
	switch k := rtype.Kind(); {
	case k == r.Bool:
		ret, okay = boolEval, true
	case numberKind(k):
		ret, okay = numEval, true
	case k == r.String:
		ret, okay = textEval, true
	case k == r.Ptr:
		ret, okay = objEval, true
	case k == r.Array || k == r.Slice:
		switch k := rtype.Elem().Kind(); {
		case numberKind(k):
			ret, okay = numListEval, true
		case k == r.String:
			ret, okay = textListEval, true
		case k == r.Ptr:
			ret, okay = objListEval, true
		}
	}
	return
}

func numberKind(k r.Kind) (ret bool) {
	switch k {
	case r.Int, r.Int8, r.Int16, r.Int32, r.Int64, r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64, r.Float32, r.Float64:
		ret = true
	}
	return
}

// switch doesnt seem to work well dstValue.Interface().(type) b/c dst is usually nil.
var boolEval = r.TypeOf((*rt.BoolEval)(nil)).Elem()
var numEval = r.TypeOf((*rt.NumberEval)(nil)).Elem()
var textEval = r.TypeOf((*rt.TextEval)(nil)).Elem()
var objEval = r.TypeOf((*rt.ObjectEval)(nil)).Elem()
var numListEval = r.TypeOf((*rt.NumListEval)(nil)).Elem()
var textListEval = r.TypeOf((*rt.TextListEval)(nil)).Elem()
var objListEval = r.TypeOf((*rt.ObjListEval)(nil)).Elem()

func arrayKind(rtype r.Type) (ret r.Kind, isArray bool) {
	if k := rtype.Kind(); k != r.Slice {
		ret = k
	} else {
		isArray = true
		ret = rtype.Elem().Kind()
	}
	return
}

// unpack the passed value
func (s *_ShadowSlot) unpack(run rt.Runtime) (ret interface{}, err error) {
	// note: we cant s.rvalue.Interface()
	// a single command can implement multiple interfaces;
	// the type switch will match whichever is listed first.
	switch rtype, val := s.rtype, s.rvalue.Interface(); rtype {
	default:
		err = errutil.New("unknown unpack type", rtype)
	case boolEval:
		if eval, ok := val.(rt.BoolEval); !ok {
			err = errutil.New("mismatched slot", rtype, s.rvalue.Type())
		} else {
			ret, err = eval.GetBool(run)
		}
	case numEval:
		if eval, ok := val.(rt.NumberEval); !ok {
			err = errutil.New("mismatched slot", rtype, s.rvalue.Type())
		} else {
			ret, err = eval.GetNumber(run)
		}
	case textEval:
		if eval, ok := val.(rt.TextEval); !ok {
			err = errutil.New("mismatched slot", rtype, s.rvalue.Type())
		} else {
			ret, err = eval.GetText(run)
		}
	case objEval:
		if eval, ok := val.(rt.ObjectEval); !ok {
			err = errutil.New("mismatched slot", rtype, s.rvalue.Type())
		} else {
			ret, err = eval.GetObject(run)
		}
	case numListEval:
		if eval, ok := val.(rt.NumListEval); !ok {
			err = errutil.New("mismatched slot", rtype, s.rvalue.Type())
		} else {
			var vals []float64
			if stream, e := eval.GetNumberStream(run); e != nil {
				err = e
			} else {
				for stream.HasNext() {
					if v, e := stream.GetNext(); e != nil {
						err = e
						break
					} else {
						vals = append(vals, v)
					}
				}
				if err == nil {
					ret = vals
				}
			}
		}
	case textListEval:
		if eval, ok := val.(rt.TextListEval); !ok {
			err = errutil.New("mismatched slot", rtype, s.rvalue.Type())
		} else {
			var vals []string
			if stream, e := eval.GetTextStream(run); e != nil {
				err = e
			} else {
				for stream.HasNext() {
					if v, e := stream.GetNext(); e != nil {
						err = e
						break
					} else {
						vals = append(vals, v)
					}
				}
				if err == nil {
					ret = vals
				}
			}
		}
	case objListEval:
		if eval, ok := val.(rt.ObjListEval); !ok {
			err = errutil.New("mismatched slot", rtype, s.rvalue.Type())
		} else {
			var vals []rt.Object
			if stream, e := eval.GetObjectStream(run); e != nil {
				err = e
			} else {
				for stream.HasNext() {
					if v, e := stream.GetNext(); e != nil {
						err = e
						break
					} else {
						vals = append(vals, v)
					}
				}
				if err == nil {
					ret = vals
				}
			}
		}
	}
	return
}
