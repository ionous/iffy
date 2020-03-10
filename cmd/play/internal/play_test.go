package internal

import (
	"testing"

	"github.com/ionous/iffy/dl/check"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rtm"
)

func TestCheck(t *testing.T) {
	prog := check.Test{TestName: "hello, goodbye",
		Go: []rt.Execute{
			&core.Choose{
				If: &core.BoolValue{Bool: true},
				True: []rt.Execute{
					&core.Say{
						Text: &core.TextValue{Text: "hello"},
					},
				},
				False: []rt.Execute{
					&core.Say{
						Text: &core.TextValue{Text: "goodbye"},
					},
				},
			},
		},
		Lines: "hello",
	}

	//run rt.Runtime
	// run, e := rtm.New(classes).Rtm()
	if run, e := rtm.New(nil).Rtm(); e != nil {
		t.Fatal(e)
	} else if e := prog.Execute(run); e != nil {
		t.Fatal(e)
	}
}
