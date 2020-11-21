package core

import (
	"testing"

	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
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
		if got, e := safe.GetText(&run, seq); e != nil {
			t.Fatal(e)
		} else if got := got.String(); got != wanted {
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

func (m *seqTest) GetField(target, field string) (ret g.Value, err error) {
	if target != object.Counter {
		err = g.UnknownField{target, field}
	} else {
		v := m.counters[field]
		ret = g.IntOf(v)
	}
	return
}

func (m *seqTest) SetField(target, field string, value g.Value) (err error) {
	if target != object.Counter {
		err = g.UnknownField{target, field}
	} else {
		m.counters[field] = value.Int()
	}
	return
}
