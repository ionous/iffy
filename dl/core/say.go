package core

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

// Buffer collects text said by other statements and returns them as a string.
// Unlike Span, it does not add or alter spaces between writes.
type Buffer struct {
	Go []rt.Execute
}

// Span collects text printed during a block and writes the text with spaces.
type Span struct {
	Go []rt.Execute
}

// Bracket sandwiches text printed during a block and puts them inside parenthesis ().
type Bracket struct {
	Go []rt.Execute
}

// Slash separates text printed during a block with left-leaning slashes.
type Slash struct {
	Go []rt.Execute
}

// Commas writes words separated with commas, ending with an "and".
type Commas struct {
	Go []rt.Execute
}

// Compose defines a spec for the composer editor.
func (*Say) Compose() composer.Spec {
	return composer.Spec{
		Name:  "say_text",
		Group: "printing",
		Desc:  "Say: print some bit of text to the player.",
	}
}

// Execute writes text to the runtime's current writer.
func (op *Say) Execute(run rt.Runtime) (err error) {
	return rt.WriteText(run, op.Text)
}

// Compose defines a spec for the composer editor.
func (*Buffer) Compose() composer.Spec {
	return composer.Spec{
		Name:  "buffer_text",
		Group: "printing",
	}
}

func (op *Buffer) GetText(run rt.Runtime) (ret string, err error) {
	var buf bytes.Buffer
	if e := rt.WritersBlock(run, &buf, func() error {
		return rt.Run(run, op.Go)
	}); e != nil {
		err = e
	} else {
		ret = buf.String()
	}
	return
}

// Compose defines a spec for the composer editor.
func (*Span) Compose() composer.Spec {
	return composer.Spec{
		Name:  "span_text",
		Group: "printing",
		Desc:  "Span Text: Writes text with spaces between words.",
	}
}

func (op *Span) GetText(run rt.Runtime) (ret string, err error) {
	var span print.Spanner // separate punctuation with spaces
	if e := rt.WritersBlock(run, &span, func() error {
		return rt.Run(run, op.Go)
	}); e != nil {
		err = e
	} else {
		ret = span.String()
	}
	return
}

// Compose defines a spec for the composer editor.
func (*Bracket) Compose() composer.Spec {
	return composer.Spec{
		Name:  "bracket_text",
		Group: "printing",
		Desc:  "Bracket text: Sandwiches text printed during a block and puts them inside parenthesis '()'.",
	}
}

func (op *Bracket) GetText(run rt.Runtime) (ret string, err error) {
	var span print.Spanner // separate punctuation with spaces
	if e := rt.WritersBlock(run, print.Bracket(&span), func() error {
		return rt.Run(run, op.Go)
	}); e != nil {
		err = e
	} else {
		ret = span.String()
	}
	return
}

// Compose defines a spec for the composer editor.
func (*Slash) Compose() composer.Spec {
	return composer.Spec{
		Name:  "slash_text",
		Group: "printing",
		Desc:  "Slash text: Separates words with left-leaning slashes '/'.",
	}
}

func (op *Slash) GetText(run rt.Runtime) (ret string, err error) {
	var span print.Spanner // separate punctuation with spaces
	if e := rt.WritersBlock(run, print.Slash(&span), func() error {
		return rt.Run(run, op.Go)
	}); e != nil {
		err = e
	} else {
		ret = span.String()
	}
	return
}

// Compose defines a spec for the composer editor.
func (*Commas) Compose() composer.Spec {
	return composer.Spec{
		Name:  "comma_text",
		Group: "printing",
		Desc:  "List text: Separates words with commas, and 'and'.",
	}
}

func (op *Commas) GetText(run rt.Runtime) (ret string, err error) {
	var span print.Spanner // separate punctuation with spaces
	if e := rt.WritersBlock(run, print.AndSeparator(&span), func() error {
		return rt.Run(run, op.Go)
	}); e != nil {
		err = e
	} else {
		ret = span.String()
	}
	return
}
