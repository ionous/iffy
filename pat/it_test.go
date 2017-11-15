package pat

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
	"github.com/ionous/sliceOf"
	"testing"
)

func TestIteration(t *testing.T) {
	// tested in reverse order
	ps := []TextListRule{
		{ListRule{Flags: Infix}, Text("d")},
		{ListRule{Flags: Postfix}, Text("c")},
		{ListRule{Flags: Prefix}, Text("b")},
		{ListRule{Flags: Postfix}, Text("a")},
	}
	if pre, post, e := split(nil, ps); e != nil {
		t.Fatal(e)
	} else if cnt := len(pre); cnt != 2 {
		t.Fatal("pre mismatched", cnt)
	} else if cnt := len(post); cnt != 2 {
		t.Fatal("post mismatched", cnt)
	} else if q, e := splitQuery(nil, ps); e != nil {
		t.Fatal(e)
	} else {
		const expected = "bdca"
		q := adaptText(nil, q)
		t.Run("direct iteration", func(t *testing.T) {
			var str string
			next := q.Iterate()
			for item, ok := next(); ok; item, ok = next() {
				n := item.(stream.ValueError)
				if v, e := n.Value, n.Error; e != nil {
					t.Fatal(e)
				} else if s := v.(string); len(s) == 0 {
					t.Fatal("empty return")
				} else if len(str) > 10 {
					t.Fatal("too many values", str)
				} else {
					str += s
				}
			}
			if str != expected {
				t.Fatal(str)
			}
		})
		t.Run("stream iteration", func(t *testing.T) {
			var str string
			for it := stream.NewTextStream(q.Iterate()); it.HasNext(); {
				if s, e := it.GetText(); e != nil {
					t.Fatal(e)
				} else if len(s) == 0 {
					t.Fatal("empty return")
				} else if len(str) > 10 {
					t.Fatal("too many values", str)
				} else {
					str += s
				}
			}
			if str != expected {
				t.Fatal(str)
			}
		})
	}
}

type Text string

func (t Text) GetTextStream(rt.Runtime) (rt.TextStream, error) {
	return stream.NewTextStream(stream.FromList(sliceOf.String(string(t)))), nil
}
