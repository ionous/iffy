package scope

import (
	"testing"

	"github.com/ionous/iffy/object"
	g "github.com/ionous/iffy/rt/generic"
)

func TestLoop(t *testing.T) {
	caps := [][]struct{ first, last bool }{
		{{true, true}},
		{{true, false}, {false, true}},
		{{true, false}, {false, false}, {false, true}},
	}
	for c, cap := range caps {
		var lf LoopFactory
		for i, cnt := 0, len(cap); i < cnt; i++ {
			cap := cap[i]
			count := i + 1
			atEnd := count == cnt
			s := lf.NextScope(g.Nothing{}, !atEnd)

			if p, e := s.GetField(object.Variables, "index"); e != nil {
				t.Fatal("loop", i, e)
			} else if fidx := p.Float(); fidx != float64(count) {
				t.Fatal("index error", fidx, "at", c, i, e)
			} else if fidx != float64(count) {
				t.Fatal("loop", i, fidx, "!=", count)
			} else if p, e := s.GetField(object.Variables, "first"); e != nil {
				t.Fatal(e)
			} else if first := p.Bool(); first != cap.first {
				t.Fatal("first error", first, "at", c, i, e)
			} else if p, e := s.GetField(object.Variables, "last"); e != nil {
				t.Fatal(e)
			} else if last := p.Bool(); last != cap.last {
				t.Fatal("last error", last, "at", c, i, e)
			} else {
				_, e := s.GetField(object.Variables, "nothing")
				if _, ok := e.(g.UnknownField); !ok {
					t.Fatal("expected loop error")
				} else {
					t.Log("loop", i, "of", cnt, fidx, first, last)
				}
			}
		}
	}
}
