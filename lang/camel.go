package lang

import (
	"bytes"
	"strings"
	"unicode"
)

// Camelize turns spaces, dashes, and underscores into words, capitalizing all but the first word,
// and lowercasing the rest of the string. likeThisIGuess.
func Camelize(name string) string {
	p := combineCase(name, true, true)
	return p.join()
}

// CombineCase is almost exactly like Camelize, only doesnt touch the case of the first rune of the first word.
func CombineCase(name string) string {
	p := combineCase(name, false, true)
	return p.join()
}

// Fields is similar to strings.Fields except this splits on dashes, case changes, spaces, and the like:
// the same rules as Camelize
func Fields(name string) []string {
	p := combineCase(name, false, false)
	return p.flush()
}

// fix. this horrible algorithm sure needs to change
// itd be fine i think to split first and then combine with word rules even
// possibly worth considering ditching camelCasing anyway:
// only support de-camelization into lower fields for template matching.
func combineCase(name string, changeFirst, changeAny bool) parts {
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
				if len(parts.arr) > 0 && changeAny {
					r = unicode.ToUpper(r)
				}
			}
			parts.WriteRune(r) // docs say err is always nil
			wasUpper = currUpper
			inword = letter
		}
	}
	return parts
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

func (p *parts) join() string {
	return strings.Join(p.flush(), "")
}
