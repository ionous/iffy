package lang

import (
	"bytes"
	"strings"
	"unicode"
)

// Camelize turns spaces, dashes, and underscores into words, capitalizing all but the first word,
// and lowercasing the rest of the string. likeThisIGuess.
func Camelize(name string) string {
	return combineCase(name, true)
}

// CombineCase is almost exactly like Camelize, only doesnt touch the case of the first rune of the first word.
func CombineCase(name string) string {
	return combineCase(name, false)
}

func combineCase(name string, changeFirst bool) (ret string) {
	type word int
	const (
		noword word = iota
		letter
		number
	)
	var parts parts
	inword, wasUpper := noword, false
	changeCase := changeFirst
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
		} else if unicode.IsLetter(r) || r == '#' {
			currUpper := unicode.IsUpper(r)
			// classify some common word changes
			sameWord := (inword == letter) && ((wasUpper == currUpper) || (wasUpper && !currUpper))
			// everything gets lowered
			if currUpper && changeCase {
				r = unicode.ToLower(r)
			}
			changeCase = true
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
