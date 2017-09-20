package play

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/rt"
	"io"
)

func startup(run rt.Runtime) (err error) {
	if story, ok := run.GetObject("story"); !ok {
		err = errutil.New("no story found")
	} else if e := rt.Determine(run, &std.Commence{story}); e != nil {
		err = e
	} else {
		var left, right string
		if e := story.GetValue("status left", &left); e != nil {
			err = e
		} else if e := story.GetValue("status right", &right); e != nil {
			err = e
		} else {
			io.WriteString(run, left)
			io.WriteString(run, right)
		}
	}
	return
}
