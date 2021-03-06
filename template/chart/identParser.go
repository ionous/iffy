package chart

type IdentParser struct {
	runes Runes
}

func (p *IdentParser) StateName() string {
	return "identifiers"
}

func (p *IdentParser) Reset() string {
	r := p.runes.String()
	p.runes = Runes{}
	return r
}

func (p *IdentParser) Identifier() string {
	return p.runes.String()
}

// first character of the identifier
func (p *IdentParser) NewRune(r rune) (ret State) {
	if isLetter(r) {
		ret = p.runes.Accept(r, Statement("ident head", p.body)) // loop...
	}
	return
}

// subsequent characters can be letters or numbers
// noting that fields are separated by dots "."
func (p *IdentParser) body(r rune) (ret State) {
	if isLetter(r) || isNumber(r) || isQualifier(r) {
		ret = p.runes.Accept(r, Statement("ident body", p.body)) // loop...
	}
	return
}
