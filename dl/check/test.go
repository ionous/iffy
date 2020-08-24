package check

import (
	"bytes"
	"log"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/print"
)

// Testing interface as a slot for testing commands
type Testing interface {
	RunTest(rt.Runtime) error
}

// TestOutput that running statements prints 'Lines'
type TestOutput struct {
	Lines string
	Go    []rt.Execute
}

func (*TestOutput) Compose() composer.Spec {
	return composer.Spec{
		Name:  "test_output",
		Spec:  "expects the output {lines|quote} when running: {activity%go+execute|ghost}",
		Group: "tests",
		Desc:  "Test Output: Run some statements, and expect that their output matches a specific value.",
	}
}

// RunTest returns an error on failure.
func (op *TestOutput) RunTest(run rt.Runtime) (err error) {
	var buf bytes.Buffer
	auto := run.Writer().(*print.AutoWriter)
	prev := auto.Target
	auto.Target = &buf
	//
	if e := rt.RunAll(run, op.Go); e != nil {
		err = errutil.New("encountered error:", e)
	} else if t := buf.String(); t != op.Lines {
		err = errutil.New("expected:", op.Lines, "got:", t)
	} else {
		auto.Target = prev
		log.Println("ok:", t)
	}
	return
}
