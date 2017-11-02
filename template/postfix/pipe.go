package postfix

// Pipe extends the shunting yard to combine chains of expressions.
// Each link in the chain becomes the last parameter of the next expression in the chain.
// For example: Greeting | Capitalize | Append: " World!"
// where greeting is "hello" would become "Hello World!"
type Pipe struct {
	Shunt
	prev Expression
}

// AddPipe delineates a link in a chain of functions.
func (p *Pipe) AddPipe() (err error) {
	if exp, e := p.Flush(); e != nil {
		err = e
	} else {
		p.prev = exp
	}
	return
}

// Flush returns the pipe's postfix ordered output, clearing the pipe.
func (p *Pipe) Flush() (ret Expression, err error) {
	if e := p.Shunt.AddExpression(p.prev); e != nil {
		err = e
	} else if exp, e := p.Shunt.Flush(); e != nil {
		err = e
	} else {
		ret = exp
		p.prev = nil
	}
	return
}
