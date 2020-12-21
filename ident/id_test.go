package ident

import (
	"strings"
	"testing"
)

func TestIds(t *testing.T) {
	testEqual(t, "#derived_class", NameOf("DerivedClass"))
}

func TestIdBasics(t *testing.T) {
	titleCase := NameOf("TitleCase")
	testEqual(t, "#title_case", titleCase, "TitleCase to camelCase")
	testEqual(t, "#two_words", NameOf("two words"), "two words to join")
	testEqual(t, "#word_dash", NameOf("word-dash"), "dashes split ids")
	testEqual(t, "#apostrophes", NameOf("apostrophe's"), "apostrophes vanish")
	testEqual(t, NameOf(""), "", "empty is as empty does")
	testEqual(t, NameOf("786_abc_123_def"), NameOf("786-abc 123 def"))
}

func TestIdNameOfs(t *testing.T) {
	testEqual(t, "#apples", NameOf("apples"))
	testEqual(t, "#apples", NameOf("Apples"))

	testEqual(t, "#apple_turnover", NameOf("apple turnover"))
	testEqual(t, "#apple_turnover", NameOf("Apple Turnover"))
	testEqual(t, "#apple_turnover", NameOf("Apple turnover"))
	testEqual(t, "#apple_turnover", NameOf("APPLE TURNOVER"))
	testEqual(t, "#apple_turnover", NameOf("apple-turnover"))

	testEqual(t, "#pascal_case", NameOf("PascalCase"))
	testEqual(t, "#camel_case", NameOf("camelCase"))

	testEqual(t, "#something_like_this", NameOf("something-like-this"))
	testEqual(t, "#allcaps", NameOf("ALLCAPS"))

	testEqual(t, "#wha_tabout_this", NameOf("whaTAboutThis"))
}

func testEqual(t *testing.T, one, two string, extra ...string) {
	if one != two {
		t.Fatal(one, two, strings.Join(extra, " "))
	}
}

// TestRecycle to ensure ids generated from ids match.
// important for gopherjs optimizations.
func TestRecycle(t *testing.T) {
	src := []string{
		"lowercase",
		"ALLCAPS",
		"PascalCase",
		"camellCase",
		"space case",
		"em-dash",
	}
	for _, src := range src {
		id := NameOf(src)
		recycledId := NameOf(id)
		testEqual(t, id, recycledId)
	}
}
