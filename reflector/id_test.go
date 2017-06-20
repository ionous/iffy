package reflector

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIds(t *testing.T) {
	assert := assert.New(t)
	dc := MakeId("DerivedClass")
	assert.Equal("$derivedClass", dc)
}

func TestIdBasics(t *testing.T) {
	titleCase := MakeId("TitleCase")
	assert.Equal(t, "$titleCase", titleCase, "TitleCase to camelCase")
	assert.Equal(t, "$twoWords", MakeId("two words"), "two words to join")
	assert.Equal(t, "$wordDash", MakeId("word-dash"), "dashes split")
	assert.Equal(t, "$apostrophes", MakeId("apostrophe's"), "apostrophes vanish")
	assert.Empty(t, MakeId(""), "empty is as empty does")
	assert.Equal(t, MakeId("786abc123def"), MakeId("786-abc 123 def"))
}

func TestIdStrings(t *testing.T) {
	assert.EqualValues(t, "$apples", MakeId("apples"))
	assert.EqualValues(t, "$apples", MakeId("Apples"))

	assert.EqualValues(t, "$appleTurnover", MakeId("apple turnover"))
	assert.EqualValues(t, "$appleTurnover", MakeId("Apple Turnover"))
	assert.EqualValues(t, "$appleTurnover", MakeId("Apple turnover"))
	assert.EqualValues(t, "$appleTurnover", MakeId("APPLE TURNOVER"))
	assert.EqualValues(t, "$appleTurnover", MakeId("apple-turnover"))

	assert.EqualValues(t, "$pascalCase", MakeId("PascalCase"))
	assert.EqualValues(t, "$camelCase", MakeId("camelCase"))

	assert.EqualValues(t, "$somethingLikeThis", MakeId("something-like-this"))
	assert.EqualValues(t, "$allcaps", MakeId("ALLCAPS"))

	assert.EqualValues(t, "$whaTaboutThis", MakeId("whaTAboutThis"))
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
		id := MakeId(src)
		recycledId := MakeId(id)
		assert.Equal(t, id, recycledId)
	}
}
