package reflector

import (
	"bytes"
	"github.com/ionous/iffy/lang"
	"strings"
	"unicode"
)

// MakeId creates a new string id from the passed raw string.
// Dashes and spaces are treated as word separators; sequences of numbers and sequences of letters are treated as separate words.
// Ported from sashimi v1 ident.Id
func MakeId(name string) (ret string) {
	if len(name) > 0 {
		if name[0] == '$' {
			ret = name
		} else {
			ret = "$" + lang.StripArticle(_MakeId(name))
		}
	}
	return
}

func _MakeId(name string) (ret string) {
	type word int
	const (
		noword word = iota
		letter
		number
	)
	var parts parts
	inword, wasUpper := noword, false

	for _, r := range name {
		if r == '-' || r == '_' || r == '=' || unicode.IsSpace(r) {
			inword = noword
			continue
		}

		if unicode.IsDigit(r) {
			if sameWord := inword == number; !sameWord {
				parts.flush()
			}
			parts.WriteRune(r)
			wasUpper = false
			inword = number
		} else if unicode.IsLetter(r) {
			currUpper := unicode.IsUpper(r)
			// classify some common word changes
			sameWord := (inword == letter) && ((wasUpper == currUpper) || (wasUpper && !currUpper))
			if currUpper {
				r = unicode.ToLower(r)
			}
			if !sameWord {
				parts.flush()
				// hack for camelCasing.
				if len(parts.arr) > 0 {
					r = unicode.ToUpper(r)
				}
			}
			parts.WriteRune(r) // docs say err is always nil
			wasUpper = currUpper
			inword = letter
		}
	}
	return strings.Join(parts.flush(), "")
	return
}

type parts struct {
	bytes.Buffer
	arr []string
}

func (p *parts) flush() []string {
	if p.Len() > 0 {
		p.arr = append(p.arr, p.String())
		p.Reset()
	}
	return p.arr
}
