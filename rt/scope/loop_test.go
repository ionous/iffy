package scope

import "testing"

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
			s := lf.NextScope(EmptyScope{}, !atEnd)
			var fidx float64
			var first, last bool

			if e := s.GetVariable("index", &fidx); e != nil || fidx != float64(count) {
				t.Fatal("index error", fidx, "at", c, i, e)
			} else if e := s.GetVariable("first", &first); e != nil || first != cap.first {
				t.Fatal("first error", first, "at", c, i, e)
			} else if e := s.GetVariable("last", &last); e != nil || last != cap.last {
				t.Fatal("last error", last, "at", c, i, e)
			} else if _, ok := s.GetVariable("nothing", nil).(UnknownVariable); !ok {
				t.Fatal("expected loop error")
			} else {
				t.Log("loop", i, "of", cnt, fidx, first, last)
			}
		}
	}
}
