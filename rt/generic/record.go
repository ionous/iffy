package generic

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/rt"
)

type Record struct {
	Nothing
	kind   string
	fields fieldData
}

var _ rt.Value = (*Record)(nil) // ensure compatibility

type fieldData map[string]rt.Value

func NewRecord(kind string, fields map[string]rt.Value) *Record {
	return &Record{kind: kind, fields: fields}
}

func (r *Record) Affinity() affine.Affinity {
	return affine.Record
}

func (r *Record) Type() string {
	return r.kind
}

func (r *Record) GetField(field string) (ret rt.Value, err error) {
	if v, ok := r.fields[field]; !ok {
		err = rt.UnknownField{r.kind, field}
	} else {
		ret = v
	}
	return
}

func (r *Record) SetField(field string, val rt.Value) (err error) {
	if v, ok := r.fields[field]; !ok {
		err = rt.UnknownField{r.kind, field}
	} else if newv, e := CopyValue(v.Affinity(), val); e != nil {
		err = e
	} else {
		r.fields[field] = newv
	}
	return
}
