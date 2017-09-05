package ops

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/kind"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/spec"
	"github.com/ionous/iffy/spec/builder"
	r "reflect"
)

type Ops struct {
	unique.Types
	ShadowTypes *unique.Stack
	Transform
}

func NewOps(classes unique.TypeRegistry) *Ops {
	return NewOpsX(classes, DefaultXform{})
}

func NewOpsX(classes unique.TypeRegistry, xform Transform) *Ops {
	return &Ops{
		make(unique.Types),
		unique.NewStack(classes),
		xform,
	}
}

type Builder struct {
	builder.Builder
}

// NewBuilder starts creating a call tree. Always returns true.
func (ops *Ops) NewBuilder(root interface{}) (*Builder, bool) {
	spec := &_Spec{cmds: ops, target: InPlace(root)}
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

// NewSpec implements spec.SpecFactory.
func (fac *_Factory) NewSpec(name string) (ret spec.Spec, err error) {
	if rtype, ok := fac.cmds.FindType(name); ok {
		ret = &_Spec{
			cmds:   fac.cmds,
			target: NewTarget(rtype),
		}
	} else if rtype, ok := fac.cmds.ShadowTypes.FindType(name); ok {
		ret = &_Spec{
			cmds:   fac.cmds,
			target: Shadow(rtype),
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
	target Target // output object we are building
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
			err = errutil.Fmt("couldnt get value for %T position %d (%s.%s)", tgt, spec.index, tgtType, fieldName)
		} else if e := spec.setField(field, arg); e != nil {
			fieldName := tgtType.Field(spec.index).Name
			err = errutil.Fmt("position %d (%s.%s) %v", spec.index, tgtType, fieldName, e)
		} else {
			spec.index++
		}
	}
	return
}

func (spec *_Spec) Assign(key string, arg interface{}) (err error) {
	field := spec.target.FieldByName(key)
	if !field.IsValid() {
		err = errutil.Fmt("couldnt get value for '%v'", key)
	} else if e := spec.setField(field, arg); e != nil {
		err = errutil.New("field", key, e)
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

// dst is the field we are setting; src the value specified in the command script.
func (spec *_Spec) setField(dst r.Value, src interface{}) (err error) {
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
		if src, e := xform(spec.cmds, src, dst); e != nil {
			err = e
		} else if e := ref.CoerceValue(dst, src); e != nil {
			err = errutil.New("couldnt assign value", e)
		}
	}
	return
}

// helper to turn dst.Value into dst.Type and transform value
func xform(t Transform, v interface{}, dst r.Value) (ret interface{}, err error) {
	if dst.Kind() == r.Interface {
		ret, err = t.TransformValue(v, dst.Type())
	} else {
		ret = v
	}
	return
}

// determine what kind of eval can produce the passed type.
func evalFromType(rtype r.Type) (ret r.Type, okay bool) {
	switch k := rtype.Kind(); {
	case k == r.Bool:
		ret, okay = boolEval, true
	case kind.IsNumber(k):
		ret, okay = numEval, true
	case k == r.String:
		ret, okay = textEval, true
	case k == r.Ptr:
		ret, okay = objEval, true
	case k == r.Interface:
		ret, okay = objEval, true
	case k == r.Array || k == r.Slice:
		switch k := rtype.Elem().Kind(); {
		case kind.IsNumber(k):
			ret, okay = numListEval, true
		case k == r.String:
			ret, okay = textListEval, true
		case k == r.Ptr:
			ret, okay = objListEval, true
		}
	}
	return
}

func arrayKind(rtype r.Type) (ret r.Kind, isArray bool) {
	if k := rtype.Kind(); k != r.Slice {
		ret = k
	} else {
		isArray = true
		ret = rtype.Elem().Kind()
	}
	return
}
