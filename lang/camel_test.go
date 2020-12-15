package lang

import (
	"testing"
)

func TestEmptyCamelize(t *testing.T) {
	// ensure the alg doesnt panic on the empty string
	if len(Camelize("")) != 0 {
		t.Fatal("empty failure")
	}
}

func TestTrim(t *testing.T) {
	test := func(pairs ...string) {
		for i, cnt := 0, len(pairs); i < cnt; i += 2 {
			a, want := pairs[i], pairs[i+1]
			if got := Camelize(a); got != want {
				t.Log(i/2, "got", got, "want", want)
				t.Fail()
			}
		}
	}
	test(
		"apple turnover", "appleTurnover",
		"  apple turnover", "appleTurnover",
		"apple turnover  ", "appleTurnover",
	)
}

func TestCamelize(t *testing.T) {
	test := func(pairs ...string) {
		for i, cnt := 0, len(pairs); i < cnt; i += 2 {
			a, want := pairs[i], pairs[i+1]
			if got := Camelize(a); got != want {
				t.Log(i/2, "got", got, "want", want)
				t.Fail()
			}
		}
	}
	test(
		// single words
		"apples", "apples",
		"Apples", "apples",
		"APPLES", "apples",
		// multiple words
		"appleTurnover", "appleTurnover",
		"apple turnover", "appleTurnover",
		"Apple Turnover", "appleTurnover",
		"Apple turnover", "appleTurnover",
		"APPLE TURNOVER", "appleTurnover",
		"apple-turnover", "appleTurnover",
		"apple---turn---over", "appleTurnOver",
		// multi-word casing
		"WasPascalCase", "wasPascalCase",
		"wasCamelCase", "wasCamelCase",
		"something-like-this", "somethingLikeThis",
		"something_like_that", "somethingLikeThat",
		"some___thing__like_that", "someThingLikeThat",
		// rando
		"whaTAboutThis", "whaTaboutThis",
		"", "",
		"lowercase", "lowercase",
		"Class", "class",
		"MyClass", "myClass",
		"MyC", "myC",
		"HTML", "html",
		"PDFLoader", "pdfloader",
		"AString", "astring",
		"SimpleXMLParser", "simpleXmlparser",
		"vimRPCPlugin", "vimRpcplugin",
		"GL11Version", "gl11Version",
		"99Bottles", "99Bottles",
		"May5", "may5",
		"BFG9000", "bfg9000",
		"BöseÜberraschung", "böseÜberraschung",
		"Two  spaces", "twoSpaces",
		"BadUTF8\xe2\xe2\xa1", "badUtf8",
	)

	// FIX: the rule should be, start with capitals then lowercase them,
	// uppercase the change if any.
	//"PDFLoader" => "pdfLoader"
	// then, if embedded, upper case them, lowercase the change.
	// "SimpleXMLParser" => "simpleXMLparser"
	// "THEcrazyCASE" => "theCRAZYcase" ?

}
