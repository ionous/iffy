package next

import (
	"bytes"

	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/print"
)

// Say some bit of text.
type Say struct {
	Text rt.TextEval
}

// Compose defines a spec for the composer editor.
func (*Say) Compose() composer.Spec {
	return composer.Spec{
		Group: "printing",
		Desc:  "Say: print some bit of text to the player.",
	}
}

// Execute writes text to the runtime's current writer.
func (op *Say) Execute(run rt.Runtime) (err error) {
	return rt.WriteText(run, op.Text)
}

// Buffer collects text said by other statements and returns them as a string.
// Unlike Span, it does not add or alter spaces between writes.
type Buffer struct {
	Block rt.Execute
}

// Compose defines a spec for the composer editor.
func (*Buffer) Compose() composer.Spec {
	return composer.Spec{
		Group: "printing",
	}
}

func (op *Buffer) GetText(run rt.Runtime) (ret string, err error) {
	var buf bytes.Buffer
	if e := rt.WritersBlock(run, &buf, func() error {
		return op.Block.Execute(run)
	}); e != nil {
		err = e
	} else {
		ret = buf.String()
	}
	return
}

// Span collects text printed during a block and writes the text with spaces.
type Span struct {
	Block rt.Execute
}

// Compose defines a spec for the composer editor.
func (*Span) Compose() composer.Spec {
	return composer.Spec{
		Group: "printing",
	}
}

func (op *Span) GetText(run rt.Runtime) (ret string, err error) {
	var span print.Spanner // separate punctuation with spaces
	if e := rt.WritersBlock(run, &span, func() error {
		return op.Block.Execute(run)
	}); e != nil {
		err = e
	} else {
		ret = span.String()
	}
	return
}

// Bracket sandwiches text printed during a block and puts them inside parenthesis ().
type Bracket struct {
	Block rt.Execute
}

// Compose defines a spec for the composer editor.
func (*Bracket) Compose() composer.Spec {
	return composer.Spec{
		Group: "printing",
		Desc:  "Bracket: sandwiches text printed during a block and puts them inside parenthesis ().",
	}
}

func (op *Bracket) GetText(run rt.Runtime) (ret string, err error) {
	var span print.Spanner // separate punctuation with spaces
	if e := rt.WritersBlock(run, print.Bracket(&span), func() error {
		return op.Block.Execute(run)
	}); e != nil {
		err = e
	} else {
		ret = span.String()
	}
	return
}

// Slash separates text printed during a block with left-leaning slashes.
type Slash struct {
	Block rt.Execute
}

// Compose defines a spec for the composer editor.
func (*Slash) Compose() composer.Spec {
	return composer.Spec{
		Group: "printing",
	}
}

func (op *Slash) GetText(run rt.Runtime) (ret string, err error) {
	var span print.Spanner // separate punctuation with spaces
	if e := rt.WritersBlock(run, print.Slash(&span), func() error {
		return op.Block.Execute(run)
	}); e != nil {
		err = e
	} else {
		ret = span.String()
	}
	return
}

// Commas writes words separated with commas, ending with an "and".
type Commas struct {
	Block rt.Execute
}

// Compose defines a spec for the composer editor.
func (*Commas) Compose() composer.Spec {
	return composer.Spec{
		Group: "printing",
	}
}

func (op *Commas) GetText(run rt.Runtime) (ret string, err error) {
	var span print.Spanner // separate punctuation with spaces
	if e := rt.WritersBlock(run, print.AndSeparator(&span), func() error {
		return op.Block.Execute(run)
	}); e != nil {
		err = e
	} else {
		ret = span.String()
	}
	return
}
