package builder

// Param targets a key-value spec argument.
// see Builder.Param(field)
type Param struct {
	block Builder
	key   string
}

// Cmd creates a new command of the passed name for the parameter mentioned by Builder.Param(). See also Builder.Cmd(). Args can contain Mementos and literals.
func (p Param) Cmd(name string, args ...interface{}) (ret *Memento) {
	if n, e := p.block.newCmd(name, args); e != nil {
		panic(e)
	} else {
		n.key = p.key
		ret = n
	}
	return
}

// Cmds creates a new array of commands for the parameter mentioned by Builder.Param(). See also Builder.Cmds()
func (p Param) Cmds(cmds ...*Memento) (ret *Memento) {
	if n, e := p.block.newCmds(cmds); e != nil {
		panic(e)
	} else {
		n.key = p.key
		ret = n
	}
	return
}

// Val specifies a single literal value for the parameter mentioned by Builder.Param(). See also Builder.Val()
func (p Param) Val(val interface{}) (ret *Memento) {
	if n, e := p.block.newVal(val); e != nil {
		panic(e)
	} else {
		n.key = p.key
		ret = n
	}
	return
}
