package ident

import (
	"strings"
	"testing"
)

func TestIds(t *testing.T) {
	testEqual(t, "#derivedClass", NameOf("DerivedClass"))
}

func TestIdBasics(t *testing.T) {
	titleCase := NameOf("TitleCase")
	testEqual(t, "#titleCase", titleCase, "TitleCase to camelCase")
	testEqual(t, "#twoWords", NameOf("two words"), "two words to join")
	testEqual(t, "#wordDash", NameOf("word-dash"), "dashes split")
	testEqual(t, "#apostrophes", NameOf("apostrophe's"), "apostrophes vanish")
	testEqual(t, NameOf(""), "", "empty is as empty does")
	testEqual(t, NameOf("786abc123def"), NameOf("786-abc 123 def"))
}

func TestIdNameOfs(t *testing.T) {
	testEqual(t, "#apples", NameOf("apples"))
	testEqual(t, "#apples", NameOf("Apples"))

	testEqual(t, "#appleTurnover", NameOf("apple turnover"))
	testEqual(t, "#appleTurnover", NameOf("Apple Turnover"))
	testEqual(t, "#appleTurnover", NameOf("Apple turnover"))
	testEqual(t, "#appleTurnover", NameOf("APPLE TURNOVER"))
	testEqual(t, "#appleTurnover", NameOf("apple-turnover"))

	testEqual(t, "#pascalCase", NameOf("PascalCase"))
	testEqual(t, "#camelCase", NameOf("camelCase"))

	testEqual(t, "#somethingLikeThis", NameOf("something-like-this"))
	testEqual(t, "#allcaps", NameOf("ALLCAPS"))

	testEqual(t, "#whaTaboutThis", NameOf("whaTAboutThis"))
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
