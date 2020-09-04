package pattern

import (
	"math"
	"testing"

	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
	"github.com/ionous/sliceOf"
)

func TestTextIteration(t *testing.T) {
	ps := []*TextListRule{
		{ListRule{Flags: Infix}, Text("1")},
		{ListRule{Flags: Postfix}, Text("2")},
		{ListRule{Flags: Prefix}, Text("3")},
		{ListRule{Filter: Skip}, Text("0")},
		{ListRule{Flags: Postfix}, Text("4")},
	}
	if inds, e := splitText(nil, ps); e != nil {
		t.Fatal(e)
	} else if cnt := len(inds); cnt != 4 {
		t.Fatal("expected 4 matching rules")
	} else {
		const expected = "3124"
		var got string
		for _, i := range inds {
			if txt := ps[i].TextListEval.(Text); len(txt) == 0 {
				t.Fatal("empty return")
			} else {
				got += string(txt)
			}
		}
		if got != expected {
			t.Fatal("got", got)
		}
		//
		t.Run("text iteration", func(t *testing.T) {
			var str string
			pat := &TextListPattern{"textList", ps}
			chain := textIterator{pat: pat, order: inds}
			it := stream.NewTextChain(&chain)

			for i := 0; it.HasNext(); i++ {
				if i >= cnt {
					t.Fatal(stream.Exceeded)
				} else {
					var txt string
					if e := it.GetNext(&txt); e != nil {
						t.Fatal(e)
					} else {
						str += txt
					}
				}
			}
			if str != expected {
				t.Fatal(str)
			}
		})
	}
}

func TestNumIteration(t *testing.T) {
	ps := []*NumListRule{
		{ListRule{Flags: Infix}, Number(1)},
		{ListRule{Filter: Skip}, Number(88)},
		{ListRule{Flags: Postfix}, Number(2)},
		{ListRule{Flags: Prefix}, Number(3)},
		{ListRule{Flags: Postfix}, Number(4)},
	}
	if inds, e := splitNumbers(nil, ps); e != nil {
		t.Fatal(e)
	} else if cnt := len(inds); cnt != 4 {
		t.Fatal("expected 4 matching rules")
	} else {
		var fin float64
		pat := &NumListPattern{"numList", ps}
		chain := numIterator{pat: pat, order: inds}
		it := stream.NewNumberChain(&chain)
		for i := 0; it.HasNext(); i++ {
			if i >= cnt {
				t.Fatal(stream.Exceeded)
			} else {
				var num float64
				if e := it.GetNext(&num); e != nil {
					t.Fatal(e)
				} else {
					fin += num * math.Pow10(cnt-i-1)
				}
			}
		}
		if fin != 3124 {
			t.Fatal("mismatched", fin)
		}
	}
}

type Text string

func (t Text) GetTextStream(rt.Runtime) (rt.Iterator, error) {
	v := string(t)
	return stream.NewTextList(sliceOf.String(v)), nil
}

type Number float64

func (n Number) GetNumberStream(rt.Runtime) (rt.Iterator, error) {
	v := float64(n)
	return stream.NewNumberList(sliceOf.Float64(v)), nil
}

type Bool bool

func (b Bool) GetBool(rt.Runtime) (bool, error) {
	return bool(b), nil
}

var Skip = Bool(false)
