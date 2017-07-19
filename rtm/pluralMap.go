package rtm

import (
	"github.com/ionous/iffy/lang"
)

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
