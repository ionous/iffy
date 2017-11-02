package chart

import (
	"fmt"
	"github.com/ionous/iffy/template/postfix"
)

// Block of text, a directive, or error.
type Block interface {
	blockNode()
}

// Blocks gathers blocks for the parsing of template data.
type Blocks struct {
	list []Block
}

// Len of the blocks accumulated thus far.
func (blocks Blocks) Len() int {
	return len(blocks.list)
}

// Blocks returns the underlying slice of accumulated blocks.
func (blocks Blocks) Blocks() []Block {
	return blocks.list
}

// AddBlock to the output.
func (blocks *Blocks) AddBlock(b Block) {
	blocks.list = append(blocks.list, b)
}

// TextBlock contains the uninterpreted parts of a template.
type TextBlock struct{ Text string }

// Directive contains the parsed content inside a directive.
// Both or either of the key and the expression can be empty.
type Directive struct {
	Key string
	postfix.Expression
}

func (b *TextBlock) String() string {
	return fmt.Sprintf("`%s`", b.Text)
}

func (*TextBlock) blockNode() {}
func (*Directive) blockNode() {}
