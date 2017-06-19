package ops

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/reflector"
	"github.com/ionous/iffy/spec"
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
		err = errutil.New("block should be a (nil) pointer (to a struct).")
	} else if structType := blockType.Elem(); structType.Kind() != r.Struct {
		err = errutil.New("block should point to a struct.")
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
		err = errutil.New("should be a pointer (to struct).")
	} else if rtype := ptrType.Elem(); rtype.Kind() != r.Struct {
		err = errutil.New("should point to a struct.")
	} else {
		id := reflector.MakeId(rtype.Name())
		if was, exists := ops.names[id]; exists && was != rtype {
			err = errutil.New("has conflicting names, id:", id, "was:", was, "type:", rtype)
		} else {
			ops.names[id] = rtype
		}
	}
	return
}

// OpBuilder implements spec.Spec.
type OpBuilder struct {
	ops       *Ops
	targetPtr r.Value // output object we are building
	index     int
}

// OpsArrayBuilder implements spec.Specs.
type OpsArrayBuilder struct {
	ops      *Ops
	cmdArray r.Value // output array we are appending to.
}

func (ops *Ops) Build(ptr interface{}) *spec.Context {
	targetPtr := r.ValueOf(ptr)
	ob := &OpBuilder{ops: ops, targetPtr: targetPtr}
	return spec.NewContext(ops, ob)
}

// NewSpec implements spc.SpecFactory.
func (ops *Ops) NewSpec(name string) (ret spec.Spec, err error) {
	id := reflector.MakeId(name)
	if rtype, ok := ops.names[id]; !ok {
		err = errutil.New("unknown command", name)
	} else {
		targetPtr := r.New(rtype)
		ret = &OpBuilder{
			ops:       ops,
			targetPtr: targetPtr,
		}
	}
	return
}

// NewSpecs implements spec.SpecFactory.
// the spec algorithm creates NewSpecs, and then assigns it to a slot
// we need the slot to targetPtr the array properly, so we just wait,
func (ops *Ops) NewSpecs() (spec.Specs, error) {
	return &OpsArrayBuilder{ops: ops}, nil
}

// Position implements Spec.
func (ob *OpBuilder) Position(arg interface{}) (err error) {
	tgt := ob.targetPtr.Elem()
	if cnt := tgt.NumField(); ob.index >= cnt {
		err = errutil.New("too many arguments. expected", ob.index)
	} else {
		field := tgt.Field(ob.index)
		if e := setField(field, arg); e != nil {
			err = errutil.New("field", ob.index, e)
		} else {
			ob.index++
		}
	}
	return
}

func (ob *OpBuilder) Assign(key string, arg interface{}) (err error) {
	id := reflector.MakeId(key)
	tgt := ob.targetPtr.Elem()
	tgtType := tgt.Type()
	for i, cnt := ob.index, tgtType.NumField(); i < cnt; i++ {
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

func (cbs *OpsArrayBuilder) AddElement(el spec.Spec) (err error) {
	if ob, ok := el.(*OpBuilder); !ok {
		err = errutil.Fmt("unexpected element type %T", el)
	} else {
		from := ob.targetPtr.Type()
		to := cbs.cmdArray.Type().Elem()
		//
		if !from.AssignableTo(to) {
			err = errutil.Fmt("incompatible element type. from: %v to: %v", from, to)
		} else {
			slice := r.Append(cbs.cmdArray, ob.targetPtr)
			cbs.cmdArray.Set(slice)
		}
	}
	return
}

// dst is field
func setField(dst r.Value, value interface{}) (err error) {
	switch src := value.(type) {
	case *OpBuilder:
		val := src.targetPtr.Interface()
		err = reflector.CoerceToValue(dst, val)
	case *OpsArrayBuilder:
		if kind, isArray := arrayKind(dst.Type()); !isArray || kind != r.Interface {
			err = errutil.New("expected an array of commands")
		} else {
			src.cmdArray = dst
		}
		// this are literals:
	case bool, float64, string, int, []float64, []string:
		err = reflector.CoerceToValue(dst, src)
	default:
		err = errutil.Fmt("assigning unexpected type %T", value)
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
