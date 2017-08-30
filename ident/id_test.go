package ident

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIds(t *testing.T) {
	assert := assert.New(t)
	dc := NameOf("DerivedClass")
	assert.Equal("$derivedClass", dc)
}

func TestIdBasics(t *testing.T) {
	titleCase := NameOf("TitleCase")
	assert.Equal(t, "$titleCase", titleCase, "TitleCase to camelCase")
	assert.Equal(t, "$twoWords", NameOf("two words"), "two words to join")
	assert.Equal(t, "$wordDash", NameOf("word-dash"), "dashes split")
	assert.Equal(t, "$apostrophes", NameOf("apostrophe's"), "apostrophes vanish")
	assert.Empty(t, NameOf(""), "empty is as empty does")
	assert.Equal(t, NameOf("786abc123def"), NameOf("786-abc 123 def"))
}

func TestIdNameOfs(t *testing.T) {
	assert.EqualValues(t, "$apples", NameOf("apples"))
	assert.EqualValues(t, "$apples", NameOf("Apples"))

	assert.EqualValues(t, "$appleTurnover", NameOf("apple turnover"))
	assert.EqualValues(t, "$appleTurnover", NameOf("Apple Turnover"))
	assert.EqualValues(t, "$appleTurnover", NameOf("Apple turnover"))
	assert.EqualValues(t, "$appleTurnover", NameOf("APPLE TURNOVER"))
	assert.EqualValues(t, "$appleTurnover", NameOf("apple-turnover"))

	assert.EqualValues(t, "$pascalCase", NameOf("PascalCase"))
	assert.EqualValues(t, "$camelCase", NameOf("camelCase"))

	assert.EqualValues(t, "$somethingLikeThis", NameOf("something-like-this"))
	assert.EqualValues(t, "$allcaps", NameOf("ALLCAPS"))

	assert.EqualValues(t, "$whaTaboutThis", NameOf("whaTAboutThis"))
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
		assert.Equal(t, id, recycledId)
	}
}
