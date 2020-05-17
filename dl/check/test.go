package check

import (
	"bytes"
	"log"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
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
		Spec:  "expect the text {lines|quote} when running: {activity%go+execute|ghost}",
		Group: "tests",
		Desc:  "Test Output: Run some statements, and expect that their output matches a specific value.",
	}
}

// RunTest returns an error on failure.
func (op *TestOutput) RunTest(run rt.Runtime) (err error) {
	var buf bytes.Buffer
	if e := rt.WritersBlock(run, &buf, func() error {
		return rt.RunAll(run, op.Go)
	}); e != nil {
		err = errutil.New("encountered error:", e)
	} else if t := buf.String(); t != op.Lines {
		err = errutil.New("expected:", op.Lines, "got:", t)
	} else {
		log.Println("ok:", t)
	}
	return
}
