package express

import (
	"regexp"
	"strings"
)

type Token struct {
	Pos   int
	Str   string
	Plain bool
}

func (f Token) String() string {
	return f.Str
}
func (f Token) Fields() []string {
	return strings.Fields(f.Str)
}

var x *regexp.Regexp

type Template []Token

// Tokenize
func Tokenize(s string) (ret Template, okay bool) {
	if ts, cnt := tokenize(s); cnt > 0 {
		ret, okay = ts, true
	}
	return
}

func tokenize(s string) (tokens []Token, symbols int) {
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
			tokens = append(tokens, Token{
				Str:   plain,
				Pos:   start,
				Plain: true,
			})
		}
		word := s[x+1 : y-1]
		start, tokens = y, append(tokens, Token{
			Str: word,
			Pos: x,
		})
		symbols++
	}
	if cnt := len(s); cnt > start {
		plain := s[start:cnt]
		tokens = append(tokens, Token{
			Str:   plain,
			Pos:   start,
			Plain: true,
		})
	}
	return
}
