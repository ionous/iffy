package builder

type Block struct {
	parent   *Memento
	mementos []*Memento
}

func (b *Block) Memento() *Memento {
	return b.mementos[len(b.mementos)-1]
}

func (b *Block) PopMemento() *Memento {
	pop := len(b.mementos) - 1
	end := b.mementos[pop]
	b.mementos = b.mementos[:pop]
	return end
}

func (b *Block) AddMemento(n *Memento) *Memento {
	b.mementos = append(b.mementos, n)
	return n
}
