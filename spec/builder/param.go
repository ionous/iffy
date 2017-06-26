package builder

// Param targets a key-value spec argument.
// see Factory.Param(field)
type Param struct {
	src *Memento
	key string
}

// Cmd creates a new command of the passed name for the parameter mentioned by Factory.Param(). See also Factory.Cmd(). Args can contain Mementos and literals.
func (p Param) Cmd(name string, args ...interface{}) (ret *Memento) {
	if n, e := p.src.factory.newCmd(p.src, name, args); e != nil {
		panic(e)
	} else {
		n.key = p.key
		ret = n
	}
	return
}

// Cmds creates a new array of commands for the parameter mentioned by Factory.Param(). See also Factory.Cmds()
func (p Param) Cmds(cmds ...*Memento) (ret *Memento) {
	if n, e := p.src.factory.newCmds(p.src, cmds); e != nil {
		panic(e)
	} else {
		n.key = p.key
		ret = n
	}
	return
}

// Val specifies a single literal value for the parameter mentioned by Factory.Param(). See also Factory.Val()
func (p Param) Val(val interface{}) (ret *Memento) {
	if n, e := p.src.factory.newVal(p.src, val); e != nil {
		panic(e)
	} else {
		n.key = p.key
		ret = n
	}
	return
}
