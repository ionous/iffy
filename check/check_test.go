package check

import (
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/print"
	"github.com/ionous/iffy/rt/writer"
)

func TestCheck(t *testing.T) {
	prog := &CheckOutput{
		Name:   "hello",
		Expect: "hello",
		Prog: core.NewActivity(
			&core.Choose{
				If: &core.Bool{Bool: true},
				True: core.NewActivity(&core.Say{
					Text: &core.Text{"hello"},
				}),
				False: core.NewActivity(&core.Say{
					Text: &core.Text{"goodbye"},
				}),
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
	rt.Panic
}
type checkTester struct {
	baseRuntime
	writer.Sink
}

func (c *checkTester) ActivateDomain(string, bool) {
}
