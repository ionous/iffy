package next

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/printer"
	"github.com/ionous/iffy/rt"
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
func (p *Say) Execute(run rt.Runtime) (err error) {
	return rt.WriteText(run, run.Writer(), p.Text)
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

func (p *Span) GetText(run rt.Runtime) (ret string, err error) {
	var span printer.Spanner // separate punctuation with spaces
	if e := rt.WritersBlock(run, &span, func() error {
		return p.Block.Execute(run)
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

func (p *Bracket) GetText(run rt.Runtime) (ret string, err error) {
	var span printer.Spanner // separate punctuation with spaces
	if e := rt.WritersBlock(run, printer.Bracket(&span), func() error {
		return p.Block.Execute(run)
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

func (p *Slash) GetText(run rt.Runtime) (ret string, err error) {
	var span printer.Spanner // separate punctuation with spaces
	if e := rt.WritersBlock(run, printer.Slash(&span), func() error {
		return p.Block.Execute(run)
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

func (p *Commas) GetText(run rt.Runtime) (ret string, err error) {
	var span printer.Spanner // separate punctuation with spaces
	if e := rt.WritersBlock(run, printer.AndSeparator(&span), func() error {
		return p.Block.Execute(run)
	}); e != nil {
		err = e
	} else {
		ret = span.String()
	}
	return
}
