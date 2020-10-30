package qna

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
)

type qnaObject struct {
	generic.Nothing
	n  *Runner // for fields cache
	id string
}

func newObjectValue(run *Runner, v interface{}) (ret rt.Value, err error) {
	if id, ok := v.(string); !ok {
		err = errutil.New("expected id value, got %v(%T)", v, v)
	} else {
		ret = &qnaObject{n: run, id: id}
	}
	return
}

func (q *qnaObject) Affinity() affine.Affinity {
	return affine.Object //
}

func (q *qnaObject) Type() string {
	// fix: should this be "kind"?
	// for now we return "object" and records will return their individual record kind
	// note: we'll have to exclude certain names from records: basically the affinities
	return "object{}"
}

func (q *qnaObject) GetField(field string) (ret rt.Value, err error) {
	// fix temp:
	var key keyType
	switch field {
	case object.Name, object.Kind, object.Kinds:
		// sigh
		key = makeKey(field, q.id)
	default:
		key = makeKey(q.id, field)
	}
	if v, e := q.n.getField(key); e != nil {
		err = e
	} else {
		ret = v.value
	}
	return
}

func (q *qnaObject) SetField(field string, val rt.Value) (err error) {
	if len(field) == 0 {
		err = errutil.Fmt("no field specified")
	} else if writable := field[0] != object.Prefix; !writable {
		err = errutil.Fmt("can't change reserved field %q", field)
	} else {
		key := makeKey(q.id, field)
		err = q.n.setField(key, val)
	}
	return
}
