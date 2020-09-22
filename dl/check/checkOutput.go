package check

import (
	"bytes"
	"log"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/print"
)

type CheckOutput struct {
	Name, Expect string
	Test         *core.Activity
}

func (t *CheckOutput) RunTest(run rt.Runtime) (err error) {
	// capture output into bytes
	var buf bytes.Buffer
	auto := run.Writer().(*print.AutoWriter)
	prev := auto.Target
	auto.Target = &buf

	run.ActivateDomain(t.Name, true)
	//
	if e := rt.RunOne(run, t.Test); e != nil {
		err = errutil.New("Test", t.Name, "encountered error:", e)
	} else if res := buf.String(); res != t.Expect {
		err = errutil.New("Test", t.Name, "expected:", res, "got:", t.Expect)
	} else {
		log.Println("test '" + t.Name + "' got '" + res + "'")
		auto.Target = prev
	}
	run.ActivateDomain(t.Name, false)
	return
}
