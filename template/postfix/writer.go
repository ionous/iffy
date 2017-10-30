package postfix

// Writer of postfix expressions.
type Writer interface {
	Write(f []Function) (n int, err error)
}

// Buffer implements Writer using a slice.
type Buffer struct {
	buf Expression
}

// Reset empties the slice.
func (b *Buffer) Reset() {
	b.buf = nil
}

// Write appends to the slice; it never fails.
func (b *Buffer) Write(f []Function) (int, error) {
	b.buf = append(b.buf, f...)
	return len(f), nil
}

// Expression returns the curret slice; it doesnt reset the slice.
func (b Buffer) Expression() Expression {
	return b.buf
}
