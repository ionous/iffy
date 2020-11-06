package qna

import (
	"bytes"
	"encoding/gob"
	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
)

type evalValue struct {
	generic.Nothing // for default method implementations
	run             rt.Runtime
	eval            interface{}
}

func newValue(run *Runner, a affine.Affinity, v interface{}) (ret rt.Value, err error) {
	switch a {
	case affine.Bool:
		ret, err = newBoolValue(run, v)
	case affine.Number:
		ret, err = newNumValue(run, v)
	case affine.Text:
		ret, err = newTextValue(run, v)
	case affine.Object:
		ret, err = newObjectValue(run, v)
	default:
		err = errutil.Fmt("unknown affinity %q", a)
	}
	return
}

func bytesToEval(b []byte, iptr interface{}) error {
	rptr := r.ValueOf(iptr)
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	return dec.DecodeValue(rptr)
}
