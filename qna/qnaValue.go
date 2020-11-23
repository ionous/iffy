package qna

import (
	"bytes"
	"encoding/gob"
	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

// take a snapshot of a cached value
// the meaning of a snapshot changes per value type.
// ex. snapshots from evals are unique instances,
// while multiple list snaps share the same slice memory.
type snapper interface {
	Snapshot(rt.Runtime) (g.Value, error)
}

type qnaValue struct {
	affinity affine.Affinity
	snapper
}

func (q qnaValue) Affinity() affine.Affinity {
	return q.affinity
}

type staticValue struct {
	affinity affine.Affinity
	val      interface{}
}

func (f staticValue) Snapshot(run rt.Runtime) (ret g.Value, err error) {
	switch i, a := f.val, f.affinity; a {
	case affine.Bool:
		switch v := i.(type) {
		case nil:
			ret = g.False // zero value for unhandled defaults in sqlite
		case bool:
			ret = g.BoolOf(v)
		case int64:
			// sqlite, boolean values can be represented as 1/0
			ret = g.BoolOf(v == 0)
		default:
			err = errutil.Fmt("unknown %s %T", a, v)
		}

	case affine.Number:
		switch v := i.(type) {
		case nil:
			ret = g.Zero // zero value for unhandled defaults in sqlite
		case int64:
			ret = g.IntOf(int(v))
		case float64:
			ret = g.FloatOf(v)
		default:
			err = errutil.Fmt("unknown %s %T", a, v)
		}

	case affine.NumList:
		switch vs := i.(type) {
		case []float64:
			ret = g.FloatsOf(vs)
		default:
			err = errutil.Fmt("unknown %s %T", a, vs)
		}

	case affine.Text:
		switch v := i.(type) {
		case nil:
			ret = g.Empty // zero value for unhandled defaults in sqlite
		case string:
			ret = g.StringOf(v)
		default:
			err = errutil.Fmt("unknown %s %T", a, v)
		}

	case affine.TextList:
		switch vs := i.(type) {
		case []string:
			ret = g.StringsOf(vs)
		default:
			err = errutil.Fmt("unknown %s %T", a, vs)
		}

	case affine.Record:
		if v, ok := i.(*g.Record); !ok {
			err = errutil.Fmt("unknown %s %T", a, i)
		} else {
			ret = g.RecordOf(v)
		}

	// we could either disallow direct record list storage, or:
	// store the requested kind for storage.
	// case affine.RecordList:
	// 	switch vs := i.(type) {
	// 	case []*g.Record:
	// 		ret = g.RecordsOf(vs)
	// 	default:
	// 		err = errutil.Fmt("unknown %s %T", a, vs)
	// 	}

	default:
		err = errutil.New("unhandled affinity", a)
	}
	return

}

type errorValue struct{ err error }

func (f errorValue) Snapshot(run rt.Runtime) (_ g.Value, err error) {
	err = f.err
	return
}

// temp, ideally.
type patternValue struct{ store interface{} }

func (f patternValue) Snapshot(run rt.Runtime) (_ g.Value, err error) {
	err = errutil.New("pattern expected use of GetEvalByName")
	return
}

func bytesToEval(b []byte, iptr interface{}) error {
	rptr := r.ValueOf(iptr)
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	return dec.DecodeValue(rptr)
}
