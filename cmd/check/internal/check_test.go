package internal

import (
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/check"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/qna"
	"github.com/ionous/iffy/rt"
)

func TestCheck(t *testing.T) {
	prog := &check.Test{
		TestName: "hello, goodbye",
		Go: &core.Choose{
			If: &core.Bool{Bool: true},
			True: &core.Say{
				Text: &core.Text{"hello"},
			},

			False: &core.Say{
				Text: &core.Text{"goodbye"},
			},
		},
		Lines: "hello",
	}
	if e := runTest(prog); e != nil {
		t.Fatal(e)
	}
}

func runTest(prog rt.BoolEval) (err error) {
	run := qna.NewRuntime(nil)
	if ok, e := rt.GetBool(run, prog); e != nil {
		err = e
	} else if !ok {
		err = errutil.New("unexpected failure", prog)
	}
	return
}
