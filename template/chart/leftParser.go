package chart

// LeftParser reads text, ending just after the opening of a directive.
// It deals with leading trim: "...{~".
type LeftParser struct {
	text, spaces []rune
	out          string
}

func (p *LeftParser) StateName() string {
	return "left parser"
}

func (p *LeftParser) GetText() string {
	return p.out
}

// NewRune starts with the first character of a string, ends just after the opening bracket of a directive.
func (p *LeftParser) NewRune(r rune) (ret State) {
	switch {
	case isOpenBracket(r):
		ret = Statement("opening", func(r rune) (ret State) {
			if !isTrim(r) {
				p.acceptSpaces()
			} else {
				ret = Terminal //end after eating this trim char
			}
			p.out = string(p.text)
			return
		})
	case r == eof:
		p.acceptSpaces()
		p.out = string(p.text)
	case isSpace(r):
		p.spaces = append(p.spaces, r)
		ret = p // loop...
	default:
		p.acceptSpaces()
		p.text = append(p.text, r)
		ret = p // loop...
	}
	return
}

func (p *LeftParser) acceptSpaces() {
	p.text, p.spaces = append(p.text, p.spaces...), nil
}
