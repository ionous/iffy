package builder

// Param targets a key-value spec argument.
// see Builder.Param(field)
type Param struct {
	block Builder
	key   string
}

// Cmd creates a new command of the passed name for the parameter mentioned by Builder.Param(). See also Builder.Cmd(). Args can contain Mementos and literals.
func (p Param) Cmd(name string, args ...interface{}) *Memento {
	n := p.block.Cmd(name, args...)
	n.key = p.key
	return n
}

// Cmds creates a new array of commands for the parameter mentioned by Builder.Param(). See also Builder.Cmds()
func (p Param) Cmds(cmds ...*Memento) *Memento {
	n := p.block.Cmds(cmds...)
	n.key = p.key
	return n
}

// Val specifies a single literal value for the parameter mentioned by Builder.Param(). See also Builder.Val()
func (p Param) Val(val interface{}) *Memento {
	n := p.block.Val(val)
	n.key = p.key
	return n
}
