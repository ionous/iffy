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
	t.Run("parts", func(t *testing.T) {
		testEqual(t, partsFn(),
			templatize(t, partStr, nil))
	})
	t.Run("cmds", func(t *testing.T) {
		classes := make(unique.Types)
		cmds := ops.NewOps(classes)
		unique.RegisterTypes(unique.PanicTypes(cmds),
			(*core.Object)(nil),
			(*TestThe)(nil))

		testEqual(t, cmdsFn(),
			templatize(t, cmdStr, cmds))
	})
}

func templatize(t *testing.T, s string, cmds *ops.Ops) (ret rt.TextEval) {
	if res, e := Templatize(Tokenize(s), cmds); e != nil {
		t.Fatal(e)
	} else {
		ret = res
	}
	return
}

func partsFn() rt.TextEval {
	return &core.Buffer{[]rt.Execute{
		&core.PrintText{
			&core.Get{
				Obj:  &core.Object{Name: "status"},
				Prop: "score",
			},
		},
		&core.PrintText{
			&core.Text{"/"},
		},
		&core.PrintText{
			&core.Get{
				Obj:  &core.Object{Name: "story"},
				Prop: "turnCount",
			},
		},
	},
	}
}

func cmdsFn() rt.TextEval {
	return &core.Buffer{[]rt.Execute{
		&core.PrintText{
			Text: &TestThe{
				&core.Object{Name: "example"},
			},
		},
	},
	}
}
