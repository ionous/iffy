package chart

import (
	"bytes"
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

// String of the blocks accumulated thus far.
func (blocks Blocks) String() string {
	type stringer interface {
		String() string
	}
	var buf bytes.Buffer
	for _, s := range blocks.list {
		buf.WriteString(s.(stringer).String())
	}
	return buf.String()
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

func (*TextBlock) blockNode() {}
func (*Directive) blockNode() {}

func (b TextBlock) String() string {
	return b.Text
}

func (b Directive) String() (ret string) {
	if len(b.Key) > 0 {
		ret = fmt.Sprintf("{%s %s}", b.Key, b.Expression)
	} else {
		ret = fmt.Sprintf("{%s}", b.Expression)
	}
	return
}
