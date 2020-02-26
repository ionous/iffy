package express

import (
	r "reflect"
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec/ops"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
)

// TestXform to verify the use of the command string converter.
func TestXform(t *testing.T) {
	classes := make(unique.Types)
	cmds := ops.NewOps(classes)

	unique.PanicBlocks(cmds,
		(*core.Commands)(nil),
		(*Commands)(nil))

	xf := NewTransform(cmds, ident.Counters{})
	t.Run("ref", func(t *testing.T) {
		ref := "{status.score}"
		testEqual(t, partsRes(),
			templatize(t, xf, ref))
	})
	t.Run("parts", func(t *testing.T) {
		parts := "{status.score}/{story.turn}"
		testEqual(t, twoPartRes(),
			templatize(t, xf, parts))
	})
	// FIX: how to say 5+5 -- printNum? why not render? is it b/c only one section?
	t.Run("ifs", func(t *testing.T) {
		ifs := "{if x}a{else}b{end}"
		testEqual(t, ifsRes(),
			templatize(t, xf, ifs))
	})
	t.Run("shuffle", func(t *testing.T) {
		shuffle := "{cycle}a{or}b{or}c{end}"
		testEqual(t, shuffleRes(),
			templatize(t, xf, shuffle))
	})
	t.Run("num", func(t *testing.T) {
		// FIX: its ugly that references can be rendered into anything.... but numbers cant.
		num := "{13|printNum?}"
		testEqual(t, numRes(),
			templatize(t, xf, num))
	})
}

func testEqual(t *testing.T, expect, res interface{}) {
	if !testify.ObjectsAreEqualValues(expect, res) {
		t.Log(pretty.Diff(res, expect))
		t.Log("got:", pretty.Sprint(res))
		t.Log("want:", pretty.Sprint(expect))
		t.FailNow()
	}
}

func templatize(t *testing.T, xform ops.Transform, s string) (ret interface{}) {
	rtype := r.TypeOf((*rt.TextEval)(nil)).Elem()
	if r, e := xform.TransformValue(r.ValueOf(s), rtype); e != nil {
		t.Fatal(e)
	} else {
		ret = r.Interface()
	}
	return
}

func numRes() rt.TextEval {
	return &core.PrintNum{&core.NumValue{13}}
}
func shuffleRes() rt.TextEval {
	return &core.CycleText{
		Id: "$cycleCounter#1",
		Values: []rt.TextEval{
			&core.TextValue{"a"},
			&core.TextValue{"b"},
			&core.TextValue{"c"},
		},
	}
}
func ifsRes() rt.TextEval {
	return &core.ChooseText{
		If:    &GetAt{Name: "x"},
		True:  &core.TextValue{"a"},
		False: &core.TextValue{"b"},
	}
}
func partsRes() rt.TextEval {
	return &Render{
		Obj:  &GetAt{"status"},
		Prop: "score",
	}
}
func twoPartRes() rt.TextEval {
	return &core.Join{[]rt.TextEval{
		&Render{
			Obj:  &GetAt{"status"},
			Prop: "score",
		},
		&core.TextValue{"/"},
		&Render{
			Obj:  &GetAt{"story"},
			Prop: "turn",
		},
	}}
}
