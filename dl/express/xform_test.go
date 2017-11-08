package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec/ops"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
	r "reflect"
	"testing"
)

// TestXform to verify the use of the command string converter.
func TestXform(t *testing.T) {
	const (
		partStr    = "{status.score}"
		twoPartStr = "{status.score}/{story.turn}"
	)
	classes := make(unique.Types)
	cmds := ops.NewOps(classes)

	unique.PanicBlocks(cmds,
		(*std.Commands)(nil),
		(*core.Commands)(nil),
		(*Commands)(nil))

	xf := NewTransform(cmds, ident.Counters{})
	t.Run("parts", func(t *testing.T) {
		testEqual(t, partsFn(),
			templatize(t, xf, partStr))
	})
	t.Run("two parts", func(t *testing.T) {
		testEqual(t, twoPartFn(),
			templatize(t, xf, twoPartStr))
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

func templatize(t *testing.T, xform ops.Transform, s string) (ret rt.TextEval) {
	rtype := r.TypeOf((*rt.TextEval)(nil)).Elem()
	if r, e := xform.TransformValue(r.ValueOf(s), rtype); e != nil {
		t.Fatal(e)
	} else {
		ret = r.Interface().(rt.TextEval)
	}
	return
}

func partsFn() rt.TextEval {
	return &Render{
		Obj:  &GetAt{"status"},
		Prop: "score",
	}
}

func twoPartFn() rt.TextEval {
	return &core.Join{[]rt.TextEval{
		&Render{
			Obj:  &GetAt{"status"},
			Prop: "score",
		},
		&core.Text{"/"},
		&Render{
			Obj:  &GetAt{"story"},
			Prop: "turn",
		},
	}}
}
