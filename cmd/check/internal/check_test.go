package internal

import (
	"testing"

	"github.com/ionous/iffy/dl/check"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/qna"
	"github.com/ionous/iffy/rt"
)

func TestCheck(t *testing.T) {
	prog := &check.TestOutput{
		"hello", []rt.Execute{
			&core.Choose{
				If: &core.Bool{Bool: true},
				True: []rt.Execute{&core.Say{
					Text: &core.Text{"hello"},
				}},
				False: []rt.Execute{&core.Say{
					Text: &core.Text{"goodbye"},
				}},
			}},
	}
	if e := runTest(prog); e != nil {
		t.Fatal(e)
	}
}

func runTest(prog check.Testing) (err error) {
	run := qna.NewRuntime(nil)
	return prog.RunTest(run)
}
