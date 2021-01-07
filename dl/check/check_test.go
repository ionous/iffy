package check

import (
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt/print"
	"github.com/ionous/iffy/rt/writer"
	"github.com/ionous/iffy/test/testutil"
)

func TestCheck(t *testing.T) {
	prog := &CheckOutput{
		Name:   "hello",
		Expect: "hello",
		Test: core.NewActivity(
			&core.ChooseAction{
				If: &core.Bool{Bool: true},
				Do: core.MakeActivity(&core.Say{
					Text: &core.Text{"hello"},
				}),
				Else: &core.ChooseNothingElse{
					Do: core.MakeActivity(&core.Say{
						Text: &core.Text{"goodbye"},
					})},
			}),
	}
	if e := runTest(prog); e != nil {
		t.Fatal(e)
	}
}

func runTest(prog *CheckOutput) (err error) {
	var run checkTester
	run.SetWriter(print.NewAutoWriter(writer.NewStdout()))
	return prog.RunTest(&run)
}

type baseRuntime struct {
	testutil.PanicRuntime
}
type checkTester struct {
	baseRuntime
	writer.Sink
}

func (c *checkTester) ActivateDomain(string, bool) {}
