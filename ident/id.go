package ident

import (
	"github.com/ionous/iffy/lang"
	"hash/fnv"
	"io"
)

type Id struct {
	Hash uint64
	Name string
}

func IdOf(str string) Id {
	name := NameOf(str)
	w := fnv.New64a()
	io.WriteString(w, name)
	sum := w.Sum64()
	return Id{sum, name}
}

// NameOf creates a new string id from the passed raw string.
// Dashes and spaces are treated as word separators; sequences of numbers and sequences of letters are treated as separate words.
// Ported from sashimi v1 ident.Id
func NameOf(str string) (ret string) {
	if len(str) > 0 {
		if str[0] == '$' {
			ret = str
		} else {
			// FIX: consider where strip article is actually needed
			ret = "$" + lang.Camelize(lang.StripArticle(str))
		}
	}
	return
}
