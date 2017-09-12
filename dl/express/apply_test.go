package express

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec/ops"
	r "reflect"
	"testing"
)

type TestThe struct {
	Obj rt.ObjectEval
}

func (*TestThe) GetText(run rt.Runtime) (ret string, err error) {
	err = errutil.New("not implemented")
	return
}

func TestApply(t *testing.T) {
	const (
		partStr    = "{status.score}"
		twoPartStr = "{status.score}/{story.turn}"
		// cmdStr     = "{go TestThe example}"
		// ifElseStr    = "{if x}{status.score}{else}{story.turnCount}{endif}"
	)
	classes := make(unique.Types)
	cmds := ops.NewOps(classes)

	unique.PanicBlocks(cmds,
		(*std.Commands)(nil),
		(*core.Commands)(nil),
		(*Commands)(nil))

	unique.PanicTypes(cmds,
		(*TestThe)(nil))

	t.Run("parts", func(t *testing.T) {
		testEqual(t, partsFn(),
			templatize(t, partStr, cmds))
	})
	t.Run("two parts", func(t *testing.T) {
		testEqual(t, twoPartFn(),
			templatize(t, twoPartStr, cmds))
	})
	// t.Run("cmds", func(t *testing.T) {
	// testEqual(t, cmdsFn(),
	// 	templatize(t, cmdStr, cmds))
	// })
}

func templatize(t *testing.T, s string, cmds *ops.Ops) (ret rt.TextEval) {
	xf := Xform{cmds: cmds}
	rtype := r.TypeOf((*rt.TextEval)(nil)).Elem()
	if r, e := xf.TransformValue(s, rtype); e != nil {
		t.Fatal(e)
	} else {
		ret = r.(rt.TextEval)
	}
	return
}

func partsFn() rt.TextEval {
	return &Render{core.Get{
		Obj:  &GetAt{Prop: "status"},
		Prop: "score",
	}}
}

func twoPartFn() rt.TextEval {
	return &core.Buffer{
		rt.ExecuteList{
			// FIX: we should be able to "say" multiple things --
			// but we need the command array interface to allow one/many/commands more transparently
			// also, maybe say should implement both get text and execute -- buffer eveerything up in the get text version.
			&core.Say{
				&Render{core.Get{
					Obj:  &GetAt{Prop: "status"},
					Prop: "score",
				}}},
			&core.Say{
				&core.Text{"/"},
			},
			&core.Say{
				&Render{core.Get{
					Obj:  &GetAt{Prop: "story"},
					Prop: "turn",
				}},
			},
		},
	}
}

func cmdsFn() rt.TextEval {
	return &core.Buffer{[]rt.Execute{
		&core.Say{
			Text: &TestThe{
				&core.Object{Name: "example"},
			},
		},
	},
	}
}
