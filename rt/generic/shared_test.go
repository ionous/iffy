package generic

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/object"

	"github.com/ionous/iffy/rt"
)

func (n *recordTest) NewKind(name string, fields []Field) (ret *Kind) {
	k := NewKind(n, name, fields)
	n.ks = append(n.ks, k)
	return k
}

func (n *recordTest) KindByName(name string) (ret *Kind, err error) {
	var ok bool
	for _, k := range n.ks {
		if k.Name() == name {
			ret, ok = k, true
			break
		}
	}
	if !ok {
		err = errutil.New("kind not found", name)
	}
	return
}

func (n *recordTest) GetField(target, field string) (ret rt.Value, err error) {
	switch target {
	case object.Variables:
		if v, ok := n.vars[field]; !ok {
			err = rt.UnknownField{target, field}
		} else {
			ret = v
		}
	default:
		err = errutil.New("unknown field", target, field)
	}
	return
}

type recordTest struct {
	rt.Panic
	ks   []*Kind
	vars map[string]*Record
}

func newRecordAccessTest() *recordTest {
	rt := new(recordTest)
	rt.NewKind("a", []Field{
		{"x", affine.Bool, "trait"},
		{"w", affine.Bool, "trait"},
		{"y", affine.Bool, "trait"},
	})
	ks := rt.NewKind("Ks", []Field{
		{"d", affine.Number, "float64"},
		{"t", affine.Text, "string"},
		{"a", affine.Text, "aspect"},
	})
	rt.vars = map[string]*Record{
		"boop": ks.NewRecord(),
		"beep": ks.NewRecord(),
	}
	return rt
}
