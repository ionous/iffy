package text

import (
	"github.com/ionous/iffy/lang"
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

// PluralRule provides a command ( the command ) for specifying single/plural pairings.
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

// Plurals holds an all lowercase mapping of single to plural pairs.
type Plurals map[string]string

// AddPlural overrides the automatic pluralization algorithm with the specified single to plural pairing.
// Compatible with the Pluralization interface.
func (p Plurals) AddPlural(single, plural string) {
	p[single] = plural
}

// Pluralize returns the plural version of a single word via table based pairs or via automatic pluralization rules.
// Compatible with the runtime pluralize interface.
func (p Plurals) Pluralize(single string) (ret string) {
	if r, ok := p[single]; ok {
		ret = r
	} else {
		ret = lang.Pluralize(single)
	}
	return
}

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
