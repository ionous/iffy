package lang

import (
	"bytes"
	"strings"
	"unicode"
)

// Breakcase turns runs of whitespace into single underscores. It does not change casing.
// "some   BIG Example" => "some_BIG_Example".
func Breakcase(name string) string {
	var b strings.Builder
	var needBreak, canBreak bool
	for _, r := range name {
		brakes := IsBreak(r)
		if ignorable := !brakes && IsIgnorable(r); ignorable {
			// consolidate all ignorable characters...
			// eventually writing a break if we've written something of note.
			needBreak = canBreak
		} else {
			// dont write a consolidated break if we are writing an explicit break
			// ex. "  _" -> just write a single underscore, not two.
			if needBreak && !brakes {
				b.WriteRune(breaker)
			}
			b.WriteRune(r)
			canBreak = !brakes
			needBreak = false
		}
	}
	return b.String()
}

// IsBreak returns true for the set of characters which breaks words in breakcase
func IsBreak(r rune) bool {
	return r == breaker
}

// IsIgnorable returns true for the set of characters which will be consolidated by breakcase
func IsIgnorable(r rune) bool {
	return unicode.IsPunct(r) || unicode.IsSpace(r)
}

const breaker = '_'

// HasBadPunct returns true for non-breakcase punctuation and spaces.
func HasBadPunct(s string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		badMatch := !IsBreak(r) && IsIgnorable(r)
		return badMatch
	}) >= 0
}

// Fields is similar to strings.Fields except this splits on dashes, case changes, spaces, and the like:
// the same rules as Camelize
func Fields(name string) []string {
	p := combineCase(name, false, false)
	return p.flush()
}

func Underscore(name string) string {
	var b strings.Builder
	for i, el := range Fields(name) {
		if i > 0 {
			b.WriteRune('_')
		}
		b.WriteString(strings.ToLower(el))
	}
	return b.String()
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
