package story

import (
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/rt"
)

type programHook interface {
	SlotType() string
	// gob requirement: return a pointer to the interface
	CmdPtr() interface{}
	// create a "pattern rule" -- no common interface currently
	NewRule(guard rt.BoolEval) (string, interface{})
}

type boolSlot struct{ cmd rt.BoolEval }

func (b *boolSlot) SlotType() string {
	return "bool_eval"
}
func (b *boolSlot) NewRule(guard rt.BoolEval) (string, interface{}) {
	return "bool_rule", &pattern.BoolRule{guard, b.cmd}
}
func (b *boolSlot) CmdPtr() interface{} {
	return &b.cmd
}

type textSlot struct{ cmd rt.TextEval }

func (b *textSlot) SlotType() string {
	return "text_eval"
}
func (b *textSlot) NewRule(guard rt.BoolEval) (string, interface{}) {
	return "text_rule", &pattern.TextRule{guard, b.cmd}
}
func (b *textSlot) CmdPtr() interface{} {
	return &b.cmd
}

type numberSlot struct{ cmd rt.NumberEval }

func (b *numberSlot) SlotType() string {
	return "number_eval"
}
func (b *numberSlot) NewRule(guard rt.BoolEval) (string, interface{}) {
	return "number_rule", &pattern.NumberRule{guard, b.cmd}
}
func (b *numberSlot) CmdPtr() interface{} {
	return &b.cmd
}

type executeSlot struct{ cmd rt.Execute }

func (b *executeSlot) SlotType() string {
	return "execute"
}
func (b *executeSlot) NewRule(guard rt.BoolEval) (string, interface{}) {
	return "execute_rule", &pattern.ExecuteRule{pattern.ListRule{Filter: guard}, b.cmd}
}
func (b *executeSlot) CmdPtr() interface{} {
	return &b.cmd
}
