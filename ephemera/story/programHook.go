package story

import (
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/rt"
)

// generates rules based on guards
type programHook interface {
	SlotType() string
	// gob requirement: return a pointer to the interface
	CmdPtr() interface{}
	// create a "pattern rule"
	// each rule returns its own kind of value -- so there's currently no common interface
	NewRule(guard rt.BoolEval, flags pattern.Flags) (string, interface{})
}

type executeSlot struct{ cmd rt.Execute }

func (b *executeSlot) SlotType() string {
	return "execute"
}
func (b *executeSlot) NewRule(guard rt.BoolEval, flags pattern.Flags) (string, interface{}) {
	return "execute_rule", &pattern.ExecuteRule{guard, flags, b.cmd}
}
func (b *executeSlot) CmdPtr() interface{} {
	return &b.cmd
}
