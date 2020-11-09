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

type staticValue struct{ val g.Value }

func (f staticValue) Snapshot(run rt.Runtime) (ret g.Value, _ error) {
	ret = f.val
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
