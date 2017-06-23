package builder

type Param struct {
	block Builder
	key   string
}

func (p Param) Cmd(name string, args ...*Memento) *Memento {
	n := p.block.Cmd(name, args...)
	n.key = p.key
	return n
}

func (p Param) Array(vals ...*Memento) *Memento {
	n := p.block.Array(vals...)
	n.key = p.key
	return n
}

func (p Param) Val(val interface{}) *Memento {
	n := p.block.Val(val)
	n.key = p.key
	return n
}
