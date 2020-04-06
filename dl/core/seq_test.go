package core

import (
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
)

func TestSequences(t *testing.T) {
	t.Run("cycle none", func(t *testing.T) {
		matchSequence(t, []string{
			"",
		}, &CycleText{Sequence{t.Name(),
			nil,
		}})
	})
	t.Run("cycle text", func(t *testing.T) {
		matchSequence(t, []string{
			"a", "b", "c", "a", "b", "c", "a",
		}, &CycleText{Sequence{t.Name(), []rt.TextEval{
			&Text{"a"},
			&Text{"b"},
			&Text{"c"},
		}}})
	})
	t.Run("stopping", func(t *testing.T) {
		matchSequence(t, []string{
			"a", "b", "c", "c", "c", "c", "c",
		}, &StoppingText{Sequence: Sequence{
			t.Name(), []rt.TextEval{
				&Text{"a"},
				&Text{"b"},
				&Text{"c"},
			}}})
	})
	t.Run("once", func(t *testing.T) {
		matchSequence(t, []string{
			"a", "", "", "", "",
		}, &StoppingText{Sequence: Sequence{
			t.Name(), []rt.TextEval{
				&Text{"a"},
			}}})
	})
	t.Run("shuffle one", func(t *testing.T) {
		matchSequence(t, []string{
			"a", "a",
		}, &ShuffleText{Sequence: Sequence{
			t.Name(), []rt.TextEval{
				&Text{"a"},
			}}})
	})
	t.Run("shuffle", func(t *testing.T) {
		matchSequence(t, []string{
			"c", "d", "b", "e", "a", "b", "e",
		}, &ShuffleText{Sequence: Sequence{
			t.Name(), []rt.TextEval{
				&Text{"a"},
				&Text{"b"},
				&Text{"c"},
				&Text{"d"},
				&Text{"e"},
			}}})
	})
}

func matchSequence(t *testing.T, want []string, seq rt.TextEval) {
	run := seqTest{counters: make(map[string]int)}
	var have []string
	for i, wanted := range want {
		if got, e := rt.GetText(&run, seq); e != nil {
			t.Fatal(e)
		} else if got != wanted {
			t.Fatalf("error at %d wanted %q got %q", i, wanted, got)
		} else {
			have = append(have, got)
		}
	}
	t.Log(t.Name(), have)
}

type seqTest struct {
	baseRuntime
	counters map[string]int
}

func (m *seqTest) Random(inclusiveMin, exclusiveMax int) int {
	return (exclusiveMax-inclusiveMin)/2 + inclusiveMin
}

func (m *seqTest) GetField(name, field string) (ret interface{}, err error) {
	if field != object.Counter {
		err = errutil.New("unknown field", field)
	} else {
		ret = m.counters[name]
	}
	return
}

func (m *seqTest) SetField(name, field string, v interface{}) (err error) {
	if field != object.Counter {
		err = errutil.New("unknown field", field)
	} else if i, ok := v.(int); !ok {
		err = errutil.New("unknown vale", field)
	} else {
		m.counters[name] = i
	}
	return
}
