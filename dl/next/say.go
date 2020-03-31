package next

import (
	"io"

	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/printer"
	"github.com/ionous/iffy/rt"
)

// Say some bit of text.
type Say struct {
	Text rt.TextWriter
}

// Compose defines a spec for the composer editor.
func (*Say) Compose() composer.Spec {
	return composer.Spec{
		Group: "printing",
		Desc:  "Say: print some bit of text to the player.",
	}
}

// Execute writes text to the runtime's current writer.
func (p *Say) Execute(run rt.Runtime) (err error) {
	return p.Text.WriteText(run, run.Writer())
}

// Span collects text printed during a block and writes the text with spaces.
type Span struct {
	Block rt.Block
}

// Compose defines a spec for the composer editor.
func (*Span) Compose() composer.Spec {
	return composer.Spec{
		Group: "printing",
	}
}

func (p *Span) WriteText(run rt.Runtime, w io.Writer) error {
	return rt.WritersBlock(run, printer.Spacing(w), func() error {
		return p.Block.Execute(run)
	})
}

// Bracket sandwiches text printed during a block and puts them inside parenthesis ().
type Bracket struct {
	Block rt.Block
}

// Compose defines a spec for the composer editor.
func (*Bracket) Compose() composer.Spec {
	return composer.Spec{
		Group: "printing",
		Desc:  "Bracket: sandwiches text printed during a block and puts them inside parenthesis ().",
	}
}

func (p *Bracket) WriteText(run rt.Runtime, w io.Writer) (err error) {
	var span printer.Spanner // separate punctuation with spaces
	if e := rt.WritersBlock(run, printer.Bracket(&span), func() error {
		return p.Block.Execute(run)
	}); e != nil {
		err = e
	} else {
		w.Write(span.Bytes())
	}
	return
}

// Slash separates text printed during a block with left-leaning slashes.
type Slash struct {
	Block rt.Block
}

// Compose defines a spec for the composer editor.
func (*Slash) Compose() composer.Spec {
	return composer.Spec{
		Group: "printing",
	}
}

func (p *Slash) WriteText(run rt.Runtime, w io.Writer) (err error) {
	var span printer.Spanner // separate punctuation with spaces
	if e := rt.WritersBlock(run, printer.Slash(&span), func() error {
		return p.Block.Execute(run)
	}); e != nil {
		err = e
	} else {
		w.Write(span.Bytes())
	}
	return
}

// Commas writes words separated with commas, ending with an "and".
type Commas struct {
	Block rt.Block
}

// Compose defines a spec for the composer editor.
func (*Commas) Compose() composer.Spec {
	return composer.Spec{
		Group: "printing",
	}
}

func (p *Commas) WriteText(run rt.Runtime, w io.Writer) (err error) {
	var span printer.Spanner // separate punctuation with spaces
	if e := rt.WritersBlock(run, printer.AndSeparator(&span), func() error {
		return p.Block.Execute(run)
	}); e != nil {
		err = e
	} else {
		w.Write(span.Bytes())
	}
	return
}
