package generic_test

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/object"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/test/testutil"
)

func (n *recordTest) NewKind(name string, fields []g.Field) (ret *g.Kind) {
	k := g.NewKind(n, name, fields)
	n.ks = append(n.ks, k)
	return k
}

func (n *recordTest) GetKindByName(name string) (ret *g.Kind, err error) {
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

func (n *recordTest) GetField(target, field string) (ret g.Value, err error) {
	switch target {
	case object.Variables:
		if v, ok := n.vars[field]; !ok {
			err = g.UnknownField{target, field}
		} else {
			ret = g.RecordOf(v)
		}
	default:
		err = errutil.New("unknown field", target, field)
	}
	return
}

type recordTest struct {
	testutil.PanicRuntime
	ks   []*g.Kind
	vars map[string]*g.Record
}

func newRecordAccessTest() *recordTest {
	rt := new(recordTest)
	rt.NewKind("a", []g.Field{
		{"x", affine.Bool, "trait"},
		{"w", affine.Bool, "trait"},
		{"y", affine.Bool, "trait"},
	})
	ks := rt.NewKind("Ks", []g.Field{
		{"d", affine.Number, "float64"},
		{"t", affine.Text, "string"},
		{"a", affine.Text, "aspect"},
	})
	rt.vars = map[string]*g.Record{
		"boop": ks.NewRecord(),
		"beep": ks.NewRecord(),
	}
	return rt
}
