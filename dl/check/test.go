package check

import (
	"bytes"
	"log"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

type Commands struct {
	*Test
}

// Test that the output of running 'Go' statements prints 'Lines'
type Test struct {
	TestName string
	Go       rt.Execute
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

// GetBool returns true if the test succeeded, otherwise it returns an error.
func (op *Test) GetBool(run rt.Runtime) (okay bool, err error) {
	var buf bytes.Buffer
	if e := rt.WritersBlock(run, &buf, func() error {
		return op.Go.Execute(run)
	}); e != nil {
		err = errutil.New("Test", op.TestName, "encountered error:", e)
	} else if t := buf.String(); t != op.Lines {
		err = errutil.New("Test", op.TestName, "expected:", op.Lines, "got:", t)
	} else {
		log.Println("Test '"+op.TestName+"':", t)
		okay = true
	}
	return
}
