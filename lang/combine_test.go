package lang

import (
	"testing"
)

// CombineCase is almost exactly like Camelize, only doesnt touch the case of the first rune of the first word.
func CombineCase(name string) string {
	p := combineCase(name, false, true)
	return p.join()
}

func TestEmptyCombine(t *testing.T) {
	// ensure the alg doesnt panic on the empty string
	if len(CombineCase("")) != 0 {
		t.Fatal("empty failure")
	}
}

func TestCombine(t *testing.T) {
	pairs := []string{
		// single words
		"apples", "apples",
		"Apples", "Apples",
		"APPLES", "Apples",
		// multi-words,
		"appleTurnover", "appleTurnover",
		"apple turnover", "appleTurnover",
		"Apple Turnover", "AppleTurnover",
		"Apple turnover", "AppleTurnover",
		"APPLE TURNOVER", "AppleTurnover",
		"apple-turnover", "appleTurnover",
		"apple---turn---over", "appleTurnOver",
		// multi-word casing,
		"WasPascalCase", "WasPascalCase",
		"wasCamelCase", "wasCamelCase",
		"something-like-this", "somethingLikeThis",
		"something_like_that", "somethingLikeThat",
		"some___thing__like_that", "someThingLikeThat",
		// rando,
		"whaTAboutThis", "whaTaboutThis",
		"", "",
		"lowercase", "lowercase",
		"Class", "Class",
		"MyClass", "MyClass",
		"MyC", "MyC",
		"HTML", "Html",
		"PDFLoader", "Pdfloader",
		"AString", "Astring",
		"SimpleXMLParser", "SimpleXmlparser",
		"vimRPCPlugin", "vimRpcplugin",
		"GL11Version", "Gl11Version",
		"99Bottles", "99Bottles",
		"May5", "May5",
		"BFG9000", "Bfg9000",
		"BöseÜberraschung", "BöseÜberraschung",
		"Two  spaces", "TwoSpaces",
		"BadUTF8\xe2\xe2\xa1", "BadUtf8",
	}
	for i, cnt := 0, len(pairs); i < cnt; i += 2 {
		test, want := pairs[i], pairs[i+1]
		if got := CombineCase(test); got != want {
			t.Fatal("wanted", want, "got", got)
		}
	}
}
