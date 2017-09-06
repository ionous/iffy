package express

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec/ops"
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
		partStr = "{status.score}/{story.turnCount}"
		cmdStr  = "{TestThe example}"
		// noneStr      = ""
		// emptyStr     = "its {} empty"
		// nobracketStr = "no quotes"
		// escapeStr    = "its {{quoted"
		// ifElseStr    = "{if x}{status.score}{else}{story.turnCount}{endif}"
	)
	t.Run("parts", func(t *testing.T) {
		testEqual(t, partsFn(),
			templatize(t, partStr, nil))
	})
	t.Run("cmds", func(t *testing.T) {
		classes := make(unique.Types)
		cmds := ops.NewOpsX(classes, core.Xform{})
		unique.RegisterTypes(unique.PanicTypes(cmds),
			(*core.Object)(nil),
			(*TestThe)(nil))

		testEqual(t, cmdsFn(),
			templatize(t, cmdStr, cmds))
	})
}

func templatize(t *testing.T, s string, cmds *ops.Ops) (ret rt.TextEval) {
	if x, ok := Tokenize(s); !ok {
		t.Fatal("couldnt tokenize", s)
	} else if res, e := Templatize(x, cmds); e != nil {
		t.Fatal(e)
	} else {
		ret = res
	}
	return
}

func partsFn() rt.TextEval {
	return &core.Buffer{[]rt.Execute{
		&core.Say{
			&core.Get{
				Obj:  &core.GetAt{"status"},
				Prop: "score",
			},
		},
		&core.Say{
			&core.Text{"/"},
		},
		&core.Say{
			&core.Get{
				Obj:  &core.GetAt{"story"},
				Prop: "turnCount",
			},
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
