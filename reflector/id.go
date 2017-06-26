package reflector

import (
	"github.com/ionous/iffy/lang"
)

// MakeId creates a new string id from the passed raw string.
// Dashes and spaces are treated as word separators; sequences of numbers and sequences of letters are treated as separate words.
// Ported from sashimi v1 ident.Id
func MakeId(name string) (ret string) {
	if len(name) > 0 {
		if name[0] == '$' {
			ret = name
		} else {
			// FIX: consider where strip article is actually needed
			ret = "$" + lang.Camelize(lang.StripArticle(name))
		}
	}
	return
}
