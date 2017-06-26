package ops

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/reflector"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec"
	"github.com/ionous/iffy/spec/builder"
	r "reflect"
)

type Ops struct {
	names map[string]r.Type
	root  []interface{}
}

// NewOps creates a registry, calling RegisterBlock on each passed element.
func NewOps(blocks ...interface{}) *Ops {
	ops := &Ops{names: make(map[string]r.Type)}
	for _, block := range blocks {
		if e := ops.RegisterBlock(block); e != nil {
			panic(e)
		}
	}
	return ops
}

// RegisterBlock registers a structure containing pointers to commands.
func (ops *Ops) RegisterBlock(block interface{}) (err error) {
	if blockType := r.TypeOf(block); blockType.Kind() != r.Ptr {
		err = errutil.New("expected (nil) pointer (to a struct).")
	} else if structType := blockType.Elem(); structType.Kind() != r.Struct {
		err = errutil.New("expected a struct pointer.")
	} else {
		for i, cnt := 0, structType.NumField(); i < cnt; i++ {
			field := structType.Field(i)
			if e := ops.registerType(field.Type); e != nil {
				err = errutil.New(field.Name, e)
				break
			}
		}
	}
	return
}

// RegisterType registers a single pointer to a command.
func (ops *Ops) RegisterType(cmd interface{}) (err error) {
	if e := ops.registerType(r.TypeOf(cmd)); e != nil {
		err = errutil.New("command", e)
	}
	return
}

// rtype should be a struct ptr.
func (ops *Ops) registerType(cmdType r.Type) (err error) {
	if ptrType := cmdType; ptrType.Kind() != r.Ptr {
		err = errutil.New("expected (nil) pointer (to a struct).")
	} else if rtype := ptrType.Elem(); rtype.Kind() != r.Struct {
		err = errutil.New("expected a struct pointer.")
	} else {
		id := reflector.MakeId(rtype.Name())
		// println("regsiter", id)
		if was, exists := ops.names[id]; exists && was != rtype {
			err = errutil.New("has conflicting names, id:", id, "was:", was, "type:", rtype)
		} else {
			ops.names[id] = rtype
		}
	}
	return
}

type Builder struct {
	builder.Builder
}

func (ops *Ops) NewBuilder(ptr interface{}) (*Builder, bool) {
	targetPtr := r.ValueOf(ptr)
	spec := &_Spec{ops: ops, targetPtr: targetPtr}
	return &Builder{
		builder.NewBuilder(&_Factory{ops}, spec),
	}, true
}

func (u *Builder) Build() (ret interface{}, err error) {
	if res, e := u.Builder.Build(); e != nil {
		err = e
	} else if spec, ok := res.(*_Spec); !ok {
		err = errutil.Fmt("unknown error")
	} else {
		ret = spec.targetPtr.Interface()
	}
	return
}

type _Factory struct {
	ops *Ops
}

// NewSpec implements spec.SpecFactory.
func (fac *_Factory) NewSpec(name string) (ret spec.Spec, err error) {
	id := reflector.MakeId(name)
	if rtype, ok := fac.ops.names[id]; !ok {
		err = errutil.New("unknown command", name, id)
	} else {
		targetPtr := r.New(rtype)
		ret = &_Spec{
			ops:       fac.ops,
			targetPtr: targetPtr,
		}
	}
	return
}

// NewSpecs implements spec.SpecFactory.
// the spec algorithm creates NewSpecs, and then assigns it to a slot
// we need the slot to targetPtr the array properly, so we just wait,
func (fac *_Factory) NewSpecs() (spec.Specs, error) {
	return &_Specs{ops: fac.ops}, nil
}

type _Spec struct {
	ops       *Ops
	targetPtr r.Value // output object we are building
	index     int
}

func (spec *_Spec) Position(arg interface{}) (err error) {
	tgt := spec.targetPtr.Elem()
	if cnt := tgt.NumField(); spec.index >= cnt {
		err = errutil.New("too many arguments. expected", cnt)
	} else {
		field := tgt.Field(spec.index)
		if e := setField(field, arg); e != nil {
			parent := spec.targetPtr.Elem().Type().Name()
			name := tgt.Type().Field(spec.index).Name
			err = errutil.Fmt("position %d (%s.%s) %v", spec.index, parent, name, e)
		} else {
			spec.index++
		}
	}
	return
}

func (spec *_Spec) Assign(key string, arg interface{}) (err error) {
	id := reflector.MakeId(key)
	tgt := spec.targetPtr.Elem()
	tgtType := tgt.Type()
	for i, cnt := spec.index, tgtType.NumField(); i < cnt; i++ {
		fieldInfo := tgtType.Field(i)
		if id == reflector.MakeId(fieldInfo.Name) {
			field := tgt.Field(i)
			if e := setField(field, arg); e != nil {
				err = errutil.New("field", key, e)
			}
			break
		}
	}
	return
}

type _Specs struct {
	ops *Ops
	els []*_Spec
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
		if e := reflector.CoerceValue(dst, src.targetPtr); e != nil {
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
				from := spec.targetPtr.Type()
				if !from.AssignableTo(elType) {
					err = errutil.Fmt("incompatible element type. from: %v to: %v", from, elType)
					break
				} else {
					slice = r.Append(slice, spec.targetPtr)
				}
			}
			dst.Set(slice)
		}

	case bool, float64, string, int, []float64, []string:
		if dst.Kind() == r.Interface {
			if literal, ok := literally(dst.Type(), src); ok {
				src = literal
			}
		}
		if e := reflector.CoerceValue(dst, src); e != nil {
			err = errutil.New("couldnt assign primitive value", e)
		}

	default:
		err = errutil.Fmt("assigning unexpected type %T", src)
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
	case []float64:
		ret = &core.Numbers{src}
		okay = true
	case []string:
		ret = &core.Texts{src}
		okay = true
	case float64:
		ret = &core.Num{src}
		okay = true
	case string:
		// could be text or object --
		switch dstType {
		case textEval:
			ret = &core.Text{src}
			okay = true
		case objEval:
			ret = &core.Object{src}
			okay = true
		}
	}
	return
}

// switch doesnt seem to work well dstValue.Interface().(type) b/c dst is usually nil.
var textEval = r.TypeOf((*rt.TextEval)(nil)).Elem()
var objEval = r.TypeOf((*rt.ObjectEval)(nil)).Elem()

func arrayKind(rtype r.Type) (ret r.Kind, isArray bool) {
	if k := rtype.Kind(); k != r.Slice {
		ret = k
	} else {
		isArray = true
		ret = rtype.Elem().Kind()
	}
	return
}
