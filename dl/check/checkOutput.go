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
		err = errutil.Fmt("ng! %s test encountered error: %s", t.Name, e)
	} else if res := buf.String(); res != t.Expect {
		err = errutil.Fmt("ng! %s test expected: %q, got: %q", t.Name, t.Expect, res)
	} else {
		log.Printf("ok. test %s got %q", t.Name, res)
		auto.Target = prev
	}
	run.ActivateDomain(t.Name, false)
	return
}
