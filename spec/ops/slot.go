package ops

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/kind"
	r "reflect"
)

// note: nothing in the slot itself guarantees that the type and value are compatible.
// that's left up to spec/ops.
type _ShadowSlot struct {
	rtype  r.Type  // type of the slot
	rvalue r.Value // spec will .Set to this value
}

// unpack the passed value
func (s *_ShadowSlot) unpack(run rt.Runtime) (ret interface{}, err error) {
	// note: we cant s.rvalue.Interface()
	// a single command can implement multiple interfaces;
	// the type switch will match whichever is listed first.
	switch rtype, val := s.rtype, s.rvalue.Interface(); rtype {
	default:
		err = errutil.New("unknown unpack type", rtype)
	case kind.BoolEval():
		if eval, ok := val.(rt.BoolEval); !ok {
			err = errutil.New("mismatched slot", rtype, s.rvalue.Type())
		} else {
			ret, err = eval.GetBool(run)
		}
	case kind.NumberEval():
		if eval, ok := val.(rt.NumberEval); !ok {
			err = errutil.New("mismatched slot", rtype, s.rvalue.Type())
		} else {
			ret, err = eval.GetNumber(run)
		}
	case kind.TextEval():
		if eval, ok := val.(rt.TextEval); !ok {
			err = errutil.New("mismatched slot", rtype, s.rvalue.Type())
		} else {
			ret, err = eval.GetText(run)
		}
	case kind.ObjectEval():
		if eval, ok := val.(rt.ObjectEval); !ok {
			err = errutil.New("mismatched slot", rtype, s.rvalue.Type())
		} else {
			ret, err = eval.GetObject(run)
			if ret == nil {
				err = errutil.Fmt("nil object from %T", eval)
			}
		}
	case kind.NumListEval():
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
	case kind.TextListEval():
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
	case kind.ObjListEval():
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
