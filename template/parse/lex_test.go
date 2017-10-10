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
	tEnd   = mkItem(item.End, "")
	tLeft  = mkItem(item.LeftBracket, "{")
	tPipe  = mkItem(item.Filter, "|")
	tRight = mkItem(item.RightBracket, "}")
	tKey   = mkItem(item.Keyword, "keyword")
	tRef   = mkItem(item.Reference, "abc")
)

//{ run: a { bob && judy } c }

var lexTests = []lexTest{
	{"empty", "", []item.Data{tEnd}},
	{"spaces", " \t\n", []item.Data{mkItem(item.Text, " \t\n"), tEnd}},
	{"text", `now is the time`, []item.Data{mkItem(item.Text, "now is the time"), tEnd}},
	{"text with comment", "hello-{/* this is a comment */}-world", []item.Data{
		mkItem(item.Text, "hello-"),
		mkItem(item.Text, "-world"),
		tEnd,
	}},
	{"empty directive", `{}`, []item.Data{tLeft, tRight, tEnd}},
	{"reference", "{abc}", []item.Data{
		tLeft,
		tRef,
		tRight,
		tEnd,
	}},
	{"dotted", "{a.b.c}", []item.Data{
		tLeft,
		mkItem(item.Reference, "a.b.c"),
		tRight,
		tEnd,
	}},
	// function
	{"function", "{player!}", []item.Data{
		tLeft,
		mkItem(item.Function, "player"),
		tRight,
		tEnd,
	}},
	// expression
	{"expression1", "{!player}", []item.Data{
		tLeft,
		mkItem(item.Expression, "!player"),
		tRight,
		tEnd,
	}},
	// expression
	{"expression2", "{5 + 5}", []item.Data{
		tLeft,
		mkItem(item.Expression, "5 + 5"),
		tRight,
		tEnd,
	}},
	{"expression3", "{ 5+5 }", []item.Data{
		tLeft,
		mkItem(item.Expression, "5+5"),
		tRight,
		tEnd,
	}},
	// function with parameters
	{"parameters", "{beep: {5} me}", []item.Data{
		tLeft,
		mkItem(item.Function, "beep"),
		tLeft, mkItem(item.Expression, "5"), tRight,
		mkItem(item.Reference, "me"),
		tRight,
		tEnd,
	}},
	// identifier filter
	{"filter id", "{mixology32|believe: my.words}", []item.Data{
		tLeft,
		mkItem(item.Reference, "mixology32"),
		tPipe,
		mkItem(item.Function, "believe"),
		mkItem(item.Reference, "my.words"),
		tRight,
		tEnd,
	}},
	// function filter
	{"filter id", "{player?|hello!|append: {5}}", []item.Data{
		tLeft,
		mkItem(item.Function, "player"),
		tPipe,
		mkItem(item.Function, "hello"),
		tPipe,
		mkItem(item.Function, "append"),
		tLeft, mkItem(item.Expression, "5"), tRight,
		tRight,
		tEnd,
	}},

	// {"for", `{for}`, []item.Data{tLeft, tFor, tRight, tEnd}},
	// {"block", `{block "foo" .}`, []item.Data{
	// 	tLeft, tBlock, tSpace, mkItem(item.String, `"foo"`), tSpace, tDot, tRight, tEnd,
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
	// interesting:
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
		tRef,
		mkItem(item.Error, "unclosed directive"),
	}},
	{"ambiguous expression", "{filter: 5}", []item.Data{
		tLeft,
		mkItem(item.Function, "filter"),
		mkItem(item.Error, "ambiguous expression"),
	}},
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
func collect(t *testing.T, src *lexTest) (ret []item.Data) {
	l := ScanText(MakeWindow(src.input, keywords))
	var last string
	for i := 0; i < 1000; i++ {
		// note: logging externally can skip states which dont emit
		// can also add a logger or change callback to lexer for better accuracy.
		if i, ok := l.Next(); !ok {
			break
		} else {
			ret = append(ret, i)
		}
		if curr := l.State(); curr != last {
			t.Log("state", curr)
			last = curr
		}
	}
	return
	// return l.Drain(1000)
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
			items := collect(t, &test)
			if !equal(items, test.items, false) {
				t.Fatalf("got\n%v\nexpected\n%v",
					itemStrings(items),
					itemStrings(test.items))
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
