package std

import (
	"github.com/ionous/iffy/rt"
)

// Pluralizer defines a spec for generating pluralization rules.
type Pluralizer interface {
	Generate(Pluralization) error
}

// Pluralization adds pairs of single/plural words.
type Pluralization interface {
	AddPlural(single, plural string)
}

// PluralRule commands the runtime pair the specified singular text with the specified plural text.
// i.e. for use with the Pluralize command.
type PluralRule struct {
	Single, Plural string
}

// Generate implements the Pluralizer spec.
func (r *PluralRule) Generate(p Pluralization) error {
	if len(r.Single) > 0 && len(r.Plural) > 0 {
		p.AddPlural(r.Single, r.Plural)
	}
	return nil
}

// Pluralize generates plural text from the passed (presumably singular) text.
type Pluralize struct {
	Text rt.TextEval
}

func (p *Pluralize) GetText(run rt.Runtime) (ret string, err error) {
	if text, e := p.Text.GetText(run); e != nil {
		err = e
	} else {
		ret = run.Pluralize(text)
	}
	return
}
