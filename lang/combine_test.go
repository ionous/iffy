package lang

import (
	"fmt"
	"testing"
)

func TestEmptyCombine(t *testing.T) {
	// ensure the alg doesnt panic on the empty string
	if len(CombineCase("")) != 0 {
		t.Fatal("empty failure")
	}
}

func ExampleCombine() {
	print := func(s ...string) {
		for _, s := range s {
			fmt.Println(fmt.Sprintf("%#v => %#v", s, CombineCase(s)))
		}
	}
	print(
		// single words
		"apples",
		"Apples",
		"APPLES",
		// multi-words
		"appleTurnover",
		"apple turnover",
		"Apple Turnover",
		"Apple turnover",
		"APPLE TURNOVER",
		"apple-turnover",
		"apple---turn---over",
		// multi-word casing
		"WasPascalCase",
		"wasCamelCase",
		"something-like-this",
		"something_like_that",
		"some___thing__like_that",
		// rando
		"whaTAboutThis",
		//fatih
		"",
		"lowercase",
		"Class",
		"MyClass",
		"MyC",
		"HTML",
		"PDFLoader",
		"AString",
		"SimpleXMLParser",
		"vimRPCPlugin",
		"GL11Version",
		"99Bottles",
		"May5",
		"BFG9000",
		"BöseÜberraschung",
		"Two  spaces",
		"BadUTF8\xe2\xe2\xa1",
	)

	// Output:
	// "apples" => "apples"
	// "Apples" => "Apples"
	// "APPLES" => "Apples"
	// "appleTurnover" => "appleTurnover"
	// "apple turnover" => "appleTurnover"
	// "Apple Turnover" => "AppleTurnover"
	// "Apple turnover" => "AppleTurnover"
	// "APPLE TURNOVER" => "AppleTurnover"
	// "apple-turnover" => "appleTurnover"
	// "apple---turn---over" => "appleTurnOver"
	// "WasPascalCase" => "WasPascalCase"
	// "wasCamelCase" => "wasCamelCase"
	// "something-like-this" => "somethingLikeThis"
	// "something_like_that" => "somethingLikeThat"
	// "some___thing__like_that" => "someThingLikeThat"
	// "whaTAboutThis" => "whaTaboutThis"
	// "" => ""
	// "lowercase" => "lowercase"
	// "Class" => "Class"
	// "MyClass" => "MyClass"
	// "MyC" => "MyC"
	// "HTML" => "Html"
	// "PDFLoader" => "Pdfloader"
	// "AString" => "Astring"
	// "SimpleXMLParser" => "SimpleXmlparser"
	// "vimRPCPlugin" => "vimRpcplugin"
	// "GL11Version" => "Gl11Version"
	// "99Bottles" => "99Bottles"
	// "May5" => "May5"
	// "BFG9000" => "Bfg9000"
	// "BöseÜberraschung" => "BöseÜberraschung"
	// "Two  spaces" => "TwoSpaces"
	// "BadUTF8\xe2\xe2\xa1" => "BadUtf8"
}
