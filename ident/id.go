package ident

import (
	"hash/fnv"
	"io"

	"github.com/ionous/iffy/lang"
)

func None() (invalid Id) {
	return
}

type Id struct {
	Hash uint64
	Name string
}

func (id Id) IsValid() bool {
	return id.Hash != 0
}

func (id Id) String() (ret string) {
	if id.Hash != 0 {
		ret = id.Name
	} else {
		ret = "<anonyomous>"
	}
	return
}

func IdOf(str string) (ret Id) {
	if name := NameOf(str); len(name) > 0 {
		w := fnv.New64a()
		io.WriteString(w, name)
		sum := w.Sum64()
		ret = Id{sum, name}
	}
	return
}

// NameOf creates a new string id from the passed raw string.
// Dashes and spaces are treated as word separators; sequences of numbers and sequences of letters are treated as separate words.
// Ported from sashimi v1 ident.Id
func NameOf(str string) (ret string) {
	if len(str) > 0 {
		if str[0] == '#' {
			ret = str
		} else {
			ret = "#" + lang.Underscore(str)
		}
	}
	return
}
