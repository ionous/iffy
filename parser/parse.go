package parser

import (
	"strings"
)

func Parse(n Scope, m Matcher, s string) {
	w := strings.Fields(s)

	ctx := Context{Scope: n, Match: m, Words: w}
	for ctx.Match != nil {
		if !ctx.Advance() {
			break
		}
	}
}
