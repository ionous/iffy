// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parse

import (
	"fmt"
	"github.com/ionous/iffy/template/item"
	"testing"
)

type lexTest struct {
	name  string
	input string
	items []item.Data
}

func mkItem(t item.Type, v string) item.Data {
	return item.Data{Type: t, Val: v}
}

var keywords = map[string]bool{"keyword": true}

var (
	tEnd = mkItem(item.End, "")
	// tFor        = mkItem(item.Identifier, "for")
	tLeft  = mkItem(item.LeftBracket, "{")
	tPipe  = mkItem(item.Filter, "|")
	tRight = mkItem(item.RightBracket, "}")
	tKey   = mkItem(item.Keyword, "keyword")
	tRef   = mkItem(item.Reference, "abc")
	// raw         = "`" + `abc\n\t\" ` + "`"
	// rawNL       = "`now is{\n}the time`" // Contains newline inside raw quote.
)

// FIX: starting always inside isnt gojng to work well
// how do you do {a}/{b} then?
// if youre inside of { } implictly, theres no whitespace eating possible

var lexTests = []lexTest{
	{"empty", "", []item.Data{tEnd}},
	{"spaces", " \t\n", []item.Data{mkItem(item.Text, " \t\n"), tEnd}},
	{"text", `now is the time`, []item.Data{mkItem(item.Text, "now is the time"), tEnd}},
	{"text with comment", "hello-{/* this is a comment */}-world", []item.Data{
		mkItem(item.Text, "hello-"),
		mkItem(item.Text, "-world"),
		tEnd,
	}},
	// {"punctuation", "{,@% }", []item.Data{
	// 	tLeft,
	// 	mkItem(item.Char, ","),
	// 	mkItem(item.Char, "@"),
	// 	mkItem(item.Char, "%"),
	// 	tSpace,
	// 	tRight,
	// 	tEnd,
	// },
	// {"parens", "{((3))}", []item.Data{
	// 	tLeft,
	// 	tLpar,
	// 	tLpar,
	// 	mkItem(item.Number, "3"),
	// 	tRpar,
	// 	tRpar,
	// 	tRight,
	// 	tEnd,
	// }},
	{"empty directive", `{}`, []item.Data{tLeft, tRight, tEnd}},
	// {"for", `{for}`, []item.Data{tLeft, tFor, tRight, tEnd}},
	// {"block", `{block "foo" .}`, []item.Data{
	// 	tLeft, tBlock, tSpace, mkItem(item.String, `"foo"`), tSpace, tDot, tRight, tEnd,
	// }},
	// {"quote", `{"abc \n\t\" "}`, []item.Data{tLeft, tQuote, tRight, tEnd}},
	// {"raw quote", "{" + raw + "}", []item.Data{tLeft, tRawQuote, tRight, tEnd}},
	// {"raw quote with newline", "{" + rawNL + "}", []item.Data{tLeft, tRawQuoteNL, tRight, tEnd}},
	// {"numbers", "{1 02 0x14 -7.2i 1e3 +1.2e-4 4.2i 1+2i}", []item.Data{
	// 	tLeft,
	// 	mkItem(item.Number, "1"),
	// 	tSpace,
	// 	mkItem(item.Number, "02"),
	// 	tSpace,
	// 	mkItem(item.Number, "0x14"),
	// 	tSpace,
	// 	mkItem(item.Number, "-7.2i"),
	// 	tSpace,
	// 	mkItem(item.Number, "1e3"),
	// 	tSpace,
	// 	mkItem(item.Number, "+1.2e-4"),
	// 	tSpace,
	// 	mkItem(item.Number, "4.2i"),
	// 	tSpace,
	// 	mkItem(item.Complex, "1+2i"),
	// 	tRight,
	// 	tEnd,
	// }},
	// {"characters", `{'a' '\n' '\'' '\\' '\u00FF' '\xFF' '本'}`, []item.Data{
	// 	tLeft,
	// 	mkItem(item.CharConstant, `'a'`),
	// 	tSpace,
	// 	mkItem(item.CharConstant, `'\n'`),
	// 	tSpace,
	// 	mkItem(item.CharConstant, `'\''`),
	// 	tSpace,
	// 	mkItem(item.CharConstant, `'\\'`),
	// 	tSpace,
	// 	mkItem(item.CharConstant, `'\u00FF'`),
	// 	tSpace,
	// 	mkItem(item.CharConstant, `'\xFF'`),
	// 	tSpace,
	// 	mkItem(item.CharConstant, `'本'`),
	// 	tRight,
	// 	tEnd,
	// }},
	// {"bools", "{true false}", []item.Data{
	// 	tLeft,
	// 	mkItem(item.Bool, "true"),
	// 	tSpace,
	// 	mkItem(item.Bool, "false"),
	// 	tRight,
	// 	tEnd,
	// }},
	// {"dot", "{.}", []item.Data{
	// 	tLeft,
	// 	tDot,
	// 	tRight,
	// 	tEnd,
	// }},
	// {"nil", "{nil}", []item.Data{
	// 	tLeft,
	// 	mkItem(item.Nil, "nil"),
	// 	tRight,
	// 	tEnd,
	// }}1qq,
	// {"dots", "{.x . .2 .x.y.z}", []item.Data{
	// 	tLeft,
	// 	mkItem(item.Field, ".x"),
	// 	tSpace,
	// 	tDot,
	// 	tSpace,
	// 	mkItem(item.Number, ".2"),
	// 	tSpace,
	// 	mkItem(item.Field, ".x"),
	// 	mkItem(item.Field, ".y"),
	// 	mkItem(item.Field, ".z"),
	// 	tRight,
	// 	tEnd,
	// }},
	// {"keywords", "{range if else end with}", []item.Data{
	// 	tLeft,
	// 	mkItem(item.Range, "range"),
	// 	tSpace,
	// 	mkItem(item.If, "if"),
	// 	tSpace,
	// 	mkItem(item.Else, "else"),
	// 	tSpace,
	// 	mkItem(item.End, "end"),
	// 	tSpace,
	// 	mkItem(item.With, "with"),
	// 	tRight,
	// 	tEnd,
	// }},
	// {"variables", "{$c := printf $ $hello $23 $ $var.Field .Method}", []item.Data{
	// 	tLeft,
	// 	mkItem(item.Variable, "$c"),
	// 	tSpace,
	// 	mkItem(item.ColonEquals, ":="),
	// 	tSpace,
	// 	mkItem(item.Identifier, "printf"),
	// 	tSpace,
	// 	mkItem(item.Variable, "$"),
	// 	tSpace,
	// 	mkItem(item.Variable, "$hello"),
	// 	tSpace,
	// 	mkItem(item.Variable, "$23"),
	// 	tSpace,
	// 	mkItem(item.Variable, "$"),
	// 	tSpace,
	// 	mkItem(item.Variable, "$var"),
	// 	mkItem(item.Field, ".Field"),
	// 	tSpace,
	// 	mkItem(item.Field, ".Method"),
	// 	tRight,
	// 	tEnd,
	// }},
	// {"variable invocation", "{$x 23}", []item.Data{
	// 	tLeft,
	// 	mkItem(item.Variable, "$x"),
	// 	tSpace,
	// 	mkItem(item.Number, "23"),
	// 	tRight,
	// 	tEnd,
	// }},
	// {"pipeline", `intro {echo hi 1.2 |noargs|args 1 "hi"} outro`, []item.Data{
	// 	mkItem(item.Text, "intro "),
	// 	tLeft,
	// 	mkItem(item.Identifier, "echo"),
	// 	tSpace,
	// 	mkItem(item.Identifier, "hi"),
	// 	tSpace,
	// 	mkItem(item.Number, "1.2"),
	// 	tSpace,
	// 	tPipe,
	// 	mkItem(item.Identifier, "noargs"),
	// 	tPipe,
	// 	mkItem(item.Identifier, "args"),
	// 	tSpace,
	// 	mkItem(item.Number, "1"),
	// 	tSpace,
	// 	mkItem(item.String, `"hi"`),
	// 	tRight,
	// 	mkItem(item.Text, " outro"),
	// 	tEnd,
	// }},
	// {"declaration", "{$v := 3}", []item.Data{
	// 	tLeft,
	// 	mkItem(item.Variable, "$v"),
	// 	tSpace,
	// 	mkItem(item.ColonEquals, ":="),
	// 	tSpace,
	// 	mkItem(item.Number, "3"),
	// 	tRight,
	// 	tEnd,
	// }},
	// {"2 declarations", "{$v , $w := 3}", []item.Data{
	// 	tLeft,
	// 	mkItem(item.Variable, "$v"),
	// 	tSpace,
	// 	mkItem(item.Char, ","),
	// 	tSpace,
	// 	mkItem(item.Variable, "$w"),
	// 	tSpace,
	// 	mkItem(item.ColonEquals, ":="),
	// 	tSpace,
	// 	mkItem(item.Number, "3"),
	// 	tRight,
	// 	tEnd,
	// }},
	// {"field of parenthesized expression", "{(.X).Y}", []item.Data{
	// 	tLeft,
	// 	tLpar,
	// 	mkItem(item.Field, ".X"),
	// 	tRpar,
	// 	mkItem(item.Field, ".Y"),
	// 	tRight,
	// 	tEnd,
	// }},
	{"trimming spaces before and after", "hello- {~~} -world", []item.Data{
		mkItem(item.Text, "hello-"),
		tLeft,
		tRight,
		mkItem(item.Text, "-world"),
		tEnd,
	}},
	{"trimming spaces before and after comment", "hello- {~/* hello */~} -world", []item.Data{
		mkItem(item.Text, "hello-"),
		mkItem(item.Text, "-world"),
		tEnd,
	}},
	{"badchar", "#{\x01}", []item.Data{
		mkItem(item.Text, "#"),
		tLeft,
		mkItem(item.Expression, "\x01"),
		tRight,
		tEnd,
	}},
	// errors
	{"linebreak directive", "{\n}", []item.Data{
		tLeft,
		mkItem(item.Error, "unexpected line break in expression"),
	}},
	{"unclosed directive", "{abc", []item.Data{
		tLeft,
		// FIXFIXFIX: it doesnt flush the reference
		tRef,
		mkItem(item.Error, "unclosed directive"),
	}},
	// {"bad identifier", "{3k}", []item.Data{
	// 	tLeft,
	// 	mkItem(item.Error, `bad number syntax: "3k"`),
	// }},
	// {"unclosed paren", "{(3}", []item.Data{
	// 	tLeft,
	// 	tLpar,
	// 	mkItem(item.Number, "3"),
	// 	mkItem(item.Error, `unclosed left paren`),
	// },
	// {"extra right paren", "{3)}", []item.Data{
	// 	tLeft,
	// 	mkItem(item.Number, "3"),
	// 	tRpar,
	// 	mkItem(item.Error, `unexpected right paren U+0029 ')'`),
	// },

	// // Fixed bugs
	// // Many elements in an directive blew the lookahead until
	// // we made lexInsideAction not loop.
	// {"long pipeline deadlock", "{|||||}", []item.Data{
	// 	tLeft,
	// 	tPipe,
	// 	tPipe,
	// 	tPipe,
	// 	tPipe,
	// 	tPipe,
	// 	tRight,
	// 	tEnd,
	// },
	{"text with bad comment", "hello-{/*/}-world", []item.Data{
		mkItem(item.Text, "hello-"),
		mkItem(item.Error, `unclosed comment`),
	}},
	{"text with comment close separated from delim", "hello-{/* */ }-world", []item.Data{
		mkItem(item.Text, "hello-"),
		mkItem(item.Error, `comment ends without closing bracket`),
	}},
	// This one is an error that we can't catch because it breaks templates with
	// minimized JavaScript. Should have fixed it before Go 1.1.
	{"unmatched right delimiter", "hello-(.}-world", []item.Data{
		mkItem(item.Text, "hello-(.}-world"),
		tEnd,
	}},
}

// collect gathers the emitted items into a slice.
func collect(t *lexTest) (items []item.Data) {
	l := ScanText(MakeWindow(t.input, keywords))
	return l.Drain(1000)
}

func equal(i1, i2 []item.Data, checkPos bool) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].Type != i2[k].Type {
			return false
		}
		if i1[k].Val != i2[k].Val {
			return false
		}
		if checkPos && i1[k].Pos != i2[k].Pos {
			return false
		}
	}
	return true
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		t.Run(test.name, func(t *testing.T) {
			items := collect(&test)
			if !equal(items, test.items, false) {
				t.Errorf("got\n%v\nexpected\n%v",
					itemStrings(items),
					itemStrings(test.items),
				)
			}
		})
	}
}

func itemStrings(items []item.Data) (ret []string) {
	for _, i := range items {
		s := fmt.Sprintf("%v(%d):'%v'\n", i.Type, i.Pos, i.Val)
		ret = append(ret, s)
	}
	return
}

// var lexPosTests = []lexTest{
// 	{"empty", "", []item.Data{tEnd}},
// 	{"punctuation", "{,@%#}", []item.Data{
// 		{itemLeftBracket, 0, "{", 1},
// 		{itemChar, 2, ",", 1},
// 		{itemChar, 3, "@", 1},
// 		{itemChar, 4, "%", 1},
// 		{itemChar, 5, "#", 1},
// 		{itemRightBracket, 6, "}", 1},
// 		{itemEnd, 8, "", 1},
// 	}},
// 	{"sample", "0123{hello}xyz", []item.Data{
// 		{ItemText, 0, "0123", 1},
// 		{itemLeftBracket, 4, "{", 1},
// 		{itemIdentifier, 6, "hello", 1},
// 		{itemRightBracket, 11, "}", 1},
// 		{ItemText, 13, "xyz", 1},
// 		{itemEnd, 16, "", 1},
// 	}},
// }

// // The other tests don't check position, to make the test cases easier to construct.
// // This one does.
// func TestPos(t *testing.T) {
// 	for _, test := range lexPosTests {
// 		items := collect(&test, "", "")
// 		if !equal(items, test.items, true) {
// 			t.Errorf("%s: got\n\t%v\nexpected\n\t%v", test.name, items, test.items)
// 			if len(items) == len(test.items) {
// 				// Detailed print; avoid item.String() to expose the position value.
// 				for i := range items {
// 					if !equal(items[i:i+1], test.items[i:i+1], true) {
// 						i1 := items[i]
// 						i2 := test.items[i]
// 						t.Errorf("\t#%d: got {%v %d %q} expected  {%v %d %q}", i, i1.item.Type, i1.currPos, i1.val, i2.item.Type, i2.currPos, i2.val)
// 					}
// 				}
// 			}
// 		}
// 	}
// }
