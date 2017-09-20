package clog_test

import (
	"github.com/ionous/iffy/spec/clog"
	"github.com/ionous/iffy/spec/cmd"
	"github.com/ionous/sliceOf"
	"os"
)

func ExampleLogging() {
	b, _ := cmd.NewBuilder()
	c := clog.Make(os.Stdout, b)
	if c.Cmd("unit").Begin() {
		if c.Param("trials").Cmds().Begin() {
			c.Cmd("match output", sliceOf.String("a", "b", "c", "d"))
			if c.Cmd("match output", sliceOf.String("a", "b", "c", "d")).Begin() {
				if c.Cmd("for each", sliceOf.Float(1, 2, 3, 4)).Begin() {
					c.Cmd("say", c.Cmd("cycle", sliceOf.String("a", "b", "c")))
					c.End()
				}
				c.End()
			}
			c.End()
		}
		c.End()
	}

	// Output:
	// Cmd unit
	//  {
	//   Param trials
	//   Cmds
	//   {
	//    Cmd match output [a b c d]
	//    Cmd match output [a b c d]
	//    {
	//     Cmd for each [1 2 3 4]
	//     {
	//      Cmd cycle [a b c]
	//      Cmd say 1 cmd/s
	//     }
	//    }
	//   }
	//  }
}
