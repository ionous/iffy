package check

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
)

type Commands struct {
	*Test
}

// Test that the output of running
type Test struct {
	TestName string
	Go       rt.ExecuteList
	Lines    string
}

func (*Test) Compose() composer.Spec {
	return composer.Spec{
		Name:  "test",
		Spec:  "For the test {test_name:text|quote}, expect the output {lines|quote} when running: {+go|ghost}.",
		Group: "literals",
		Desc:  "Bool Value: specify an explicit true or false value.",
	}
}

// Execute evals, eats the returns
func (op *Test) Execute(run rt.Runtime) (err error) {
	b := core.Buffer{op.Go}
	if t, e := b.GetText(run); e != nil {
		err = errutil.New("Test", op.TestName, "encountered error:", e)
	} else if t != op.Lines {
		err = errutil.New("Test", op.TestName, "expected:", op.Lines, "got:", t)
	} else {
		println("test", op.TestName, ":", t)
	}
	return
}
