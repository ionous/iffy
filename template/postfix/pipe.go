package postfix

// Pipe extends the shunting yard to combine chains of expressions.
// Each link in the chain becomes the last parameter of the next expression in the chain.
// For example: Greeting | Capitalize | Append: " World!"
// where greeting is "hello" would become "Hello World!"
type Pipe struct {
	Shunt
	w    Writer
	buf  Buffer
	prev Expression
}

func (p *Pipe) Reset(w Writer) {
	p.w = w
	p.prev = nil
	p.buf.Reset()
	p.Shunt.Reset(&p.buf)
	return
}

// AddPipe delineates a link in a chain of functions.
func (p *Pipe) AddPipe() (err error) {
	if _, e := p.buf.Write(p.prev); e != nil {
		err = e
	} else if e := p.Shunt.Flush(); e != nil {
		err = e
	} else {
		p.prev = p.buf.Expression()
		p.buf.Reset()
	}
	return
}

// Flush returns the pipe's postfix ordered output, clearing the pipe.
func (p *Pipe) Flush() (err error) {
	if _, e := p.buf.Write(p.prev); e != nil {
		err = e
	} else if e := p.Shunt.Flush(); e != nil {
		err = e
	} else {
		p.w.Write(p.buf.Expression())
		p.buf.Reset()
		p.prev = nil
	}
	return
}
