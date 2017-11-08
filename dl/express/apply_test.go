package express

// import (
// 	"github.com/ionous/errutil"
// 	"github.com/ionous/iffy/dl/core"
// 	"github.com/ionous/iffy/dl/std"
// 	"github.com/ionous/iffy/ident"
// 	"github.com/ionous/iffy/ref/unique"
// 	"github.com/ionous/iffy/rt"
// 	"github.com/ionous/iffy/spec/ops"
// 	r "reflect"
// 	"testing"
// )

// type TestThe struct{ rt.ObjectEval }

// func (*TestThe) GetText(run rt.Runtime) (ret string, err error) {
// 	err = errutil.New("not implemented")
// 	return
// }

// func TestApply(t *testing.T) {
// 	const (
// 		partStr    = "{status.score}"
// 		twoPartStr = "{status.score}/{story.turn}"
// 	)
// 	classes := make(unique.Types)
// 	cmds := ops.NewOps(classes)

// 	unique.PanicBlocks(cmds,
// 		(*std.Commands)(nil),
// 		(*core.Commands)(nil),
// 		(*Commands)(nil))

// 	unique.PanicTypes(cmds,
// 		(*TestThe)(nil))

// 	xf := MakeXform(cmds, ident.Counters{})

// 	t.Run("parts", func(t *testing.T) {
// 		testEqual(t, partsFn(),
// 			templatize(t, xf, partStr))
// 	})
// 	t.Run("two parts", func(t *testing.T) {
// 		testEqual(t, twoPartFn(),
// 			templatize(t, xf, twoPartStr))
// 	})
// }

// func templatize(t *testing.T, xform ops.Transform, s string) (ret rt.TextEval) {
// 	rtype := r.TypeOf((*rt.TextEval)(nil)).Elem()
// 	if r, e := xform.TransformValue(r.ValueOf(s), rtype); e != nil {
// 		t.Fatal(e)
// 	} else {
// 		ret = r.Interface().(rt.TextEval)
// 	}
// 	return
// }

// func partsFn() rt.TextEval {
// 	return &Render{
// 		Obj:  &GetAt{"status"},
// 		Prop: "score",
// 	}
// }

// func twoPartFn() rt.TextEval {
// 	return &core.Join{[]rt.TextEval{
// 		&Render{
// 			Obj:  &GetAt{"status"},
// 			Prop: "score",
// 		},
// 		&core.Text{"/"},
// 		&Render{
// 			Obj:  &GetAt{"story"},
// 			Prop: "turn",
// 		},
// 	}}
// }
