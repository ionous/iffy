package template

import (
	"regexp"
	"strings"
)

// Token contains either a run of plain text, or a run of text from inside braces.
// Unlike a traditional token, a token's text can contain whitespaces.
type Token struct {
	Pos   int
	Str   string
	Plain bool
}

// String returns the token's text.
func (f Token) String() string {
	return f.Str
}

// Fields splits the token into whitespace chunks.
func (f Token) Fields() []string {
	return strings.Fields(f.Str)
}

// CheckFor determines if the token starts with the given text.
func (f Token) CheckFor(name string) (ret []string, okay bool) {
	if !f.Plain {
		parts := strings.Fields(f.Str)
		if len(parts) > 0 {
			if g := parts[0]; strings.EqualFold(g, name) {
				ret, okay = parts[1:], true
			}
		}
	}
	return
}

var x *regexp.Regexp

// Tokenize splits the string into tokens.
func Tokenize(s string) (ret []Token) {
	var ts []Token
	if x == nil {
		x = regexp.MustCompile(`({}|{[^{][^}]*})`)
	}
	res := x.FindAllStringIndex(s, -1)
	start := 0
	//[[0 6] [11 17] [23 30]]
	//[[3 6]]
	for _, l := range res {
		x, y := l[0], l[1]
		if x > start {
			plain := s[start:x]
			ts = append(ts, Token{
				Str:   plain,
				Pos:   start,
				Plain: true,
			})
		}
		word := s[x+1 : y-1]
		start, ts = y, append(ts, Token{
			Str: word,
			Pos: x,
		})
	}
	if cnt := len(s); cnt > start {
		plain := s[start:cnt]
		ts = append(ts, Token{
			Str:   plain,
			Pos:   start,
			Plain: true,
		})
	}
	if len(ts) > 1 || (len(ts) == 1 && !ts[0].Plain) {
		ret = ts
	}
	return
}
