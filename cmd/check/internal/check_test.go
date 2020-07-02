package internal

import (
	"testing"

	"github.com/ionous/iffy/dl/check"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/print"
	"github.com/ionous/iffy/rt/writer"
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
