package chart

type blockParser struct {
	blocks   []Block
	text     []rune
	spaces   []rune
	newBlock blockFactory
}

// creates new subBlockParser
type blockFactory func() subBlockParser

// the primary implementation is directiveParser
type subBlockParser interface {
	NewRune(rune) State
	GetBlock() (Block, error)
}

// GetBlocks doesnt return an error directly.
func (p *blockParser) GetBlocks() []Block {
	return p.blocks
}

func (p *blockParser) NewRune(r rune) (ret State) {
	switch {
	case r == eof:
		p.flushText(false)
		break // done.

	case isOpenBracket(r):
		ret = Statement(p.afterBracket)

	case isSpace(r):
		p.spaces = append(p.spaces, r)
		ret = p // loop...

	default:
		p.text = append(p.text, r)
		ret = p // loop...
	}
	return
}

func (p *blockParser) afterBracket(r rune) (ret State) {
	// write any pending text
	trim := isTrim(r)
	p.flushText(trim)
	dir := p.newBlock()
	next := makeChain(spaces, makeChain(dir, Statement(func(r rune) State {
		if d, e := dir.GetBlock(); e != nil {
			err := ErrorBlock{e}
			p.addBlock(err)
		} else if d != nil {
			p.addBlock(d)
		}
		// regardless of the error, loop:
		return p.NewRune(r)
	})))
	// if the rune was the trim char, then we can skip it;
	// otherwise, we will need to read it.
	if trim {
		ret = next
	} else {
		ret = next.NewRune(r)
	}
	return
}

// write any queued text as a block
// if trim is true, we skip trailing spaces, otherwise we write those too.
func (p *blockParser) flushText(trim bool) {
	text, spaces := p.text, p.spaces
	p.text, p.spaces = nil, nil
	if !trim {
		text = append(text, spaces...)
	}
	if len(text) > 0 {
		block := TextBlock{string(text)}
		p.addBlock(block)
	}
}

func (p *blockParser) addBlock(b Block) {
	p.blocks = append(p.blocks, b)
}
