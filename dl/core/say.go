package core

import (
	"bytes"

	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/print"
	"github.com/ionous/iffy/rt/safe"
	"github.com/ionous/iffy/rt/writer"
)

// Say some bit of text.
type Say struct {
	Text rt.TextEval
}

// Buffer collects text said by other statements and returns them as a string.
// Unlike Span, it does not add or alter spaces between writes.
type Buffer struct {
	Go *Activity
}

// Span collects text printed during a block and writes the text with spaces.
type Span struct {
	Go *Activity
}

// Bracket sandwiches text printed during a block and puts them inside parenthesis ().
type Bracket struct {
	Go *Activity
}

// Slash separates text printed during a block with left-leaning slashes.
type Slash struct {
	Go *Activity
}

// Commas writes words separated with commas, ending with an "and".
type Commas struct {
	Go *Activity
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
	return safe.WriteText(run, op.Text)
}

// Compose defines a spec for the composer editor.
func (*Buffer) Compose() composer.Spec {
	return composer.Spec{
		Name:  "buffer_text",
		Group: "printing",
	}
}

func (op *Buffer) GetText(run rt.Runtime) (g.Value, error) {
	var buf bytes.Buffer
	return writeSpan(run, &buf, op, op.Go, &buf)
}

// Compose defines a spec for the composer editor.
func (*Span) Compose() composer.Spec {
	return composer.Spec{
		Name:  "span_text",
		Group: "printing",
		Desc:  "Span Text: Writes text with spaces between words.",
	}
}

func (op *Span) GetText(run rt.Runtime) (g.Value, error) {
	span := print.NewSpanner() // separate punctuation with spaces
	return writeSpan(run, span, op, op.Go, span)
}

// Compose defines a spec for the composer editor.
func (*Bracket) Compose() composer.Spec {
	return composer.Spec{
		Name:  "bracket_text",
		Group: "printing",
		Desc:  "Bracket text: Sandwiches text printed during a block and puts them inside parenthesis '()'.",
	}
}

func (op *Bracket) GetText(run rt.Runtime) (g.Value, error) {
	span := print.NewSpanner() // separate punctuation with spaces
	return writeSpan(run, span, op, op.Go, print.Parens(span))
}

// Compose defines a spec for the composer editor.
func (*Slash) Compose() composer.Spec {
	return composer.Spec{
		Name:  "slash_text",
		Group: "printing",
		Desc:  "Slash text: Separates words with left-leaning slashes '/'.",
	}
}

func (op *Slash) GetText(run rt.Runtime) (g.Value, error) {
	span := print.NewSpanner() // separate punctuation with spaces
	return writeSpan(run, span, op, op.Go, print.Slash(span))
}

// Compose defines a spec for the composer editor.
func (*Commas) Compose() composer.Spec {
	return composer.Spec{
		Name:  "comma_text",
		Group: "printing",
		Desc:  "List text: Separates words with commas, and 'and'.",
	}
}

func (op *Commas) GetText(run rt.Runtime) (g.Value, error) {
	span := print.NewSpanner() // separate punctuation with spaces
	return writeSpan(run, span, op, op.Go, print.AndSeparator(span))
}

type stringer interface{ String() string }

func writeSpan(run rt.Runtime, span stringer, op composer.Slat, act *Activity, w writer.Output) (ret g.Value, err error) {
	if e := rt.WritersBlock(run, w, func() error {
		return safe.Run(run, act)
	}); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.StringOf(span.String())
	}
	return
}
