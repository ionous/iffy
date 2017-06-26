package lang

import "fmt"

func ExampleCamelize() {
	print := func(s ...string) {
		for _, s := range s {
			fmt.Println(fmt.Sprintf("%#v => %#v", s, Camelize(s)))
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
		// multi-word casing
		"WasPascalCase",
		"wasCamelCase",
		"something-like-this",
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

	// FIX: the rule should be, start with capitals then lowercase them,
	// uppercase the change if any.
	//"PDFLoader" => "pdfLoader"
	// then, if embedded, upper case them, lowercase the change.
	// "SimpleXMLParser" => "simpleXMLparser"
	// "THEcrazyCASE" => "theCRAZYcase" ?

	// Output:
	// "apples" => "apples"
	// "Apples" => "apples"
	// "APPLES" => "apples"
	// "appleTurnover" => "appleTurnover"
	// "apple turnover" => "appleTurnover"
	// "Apple Turnover" => "appleTurnover"
	// "Apple turnover" => "appleTurnover"
	// "APPLE TURNOVER" => "appleTurnover"
	// "apple-turnover" => "appleTurnover"
	// "WasPascalCase" => "wasPascalCase"
	// "wasCamelCase" => "wasCamelCase"
	// "something-like-this" => "somethingLikeThis"
	// "whaTAboutThis" => "whaTaboutThis"
	// "" => ""
	// "lowercase" => "lowercase"
	// "Class" => "class"
	// "MyClass" => "myClass"
	// "MyC" => "myC"
	// "HTML" => "html"
	// "PDFLoader" => "pdfloader"
	// "AString" => "astring"
	// "SimpleXMLParser" => "simpleXmlparser"
	// "vimRPCPlugin" => "vimRpcplugin"
	// "GL11Version" => "gl11Version"
	// "99Bottles" => "99Bottles"
	// "May5" => "may5"
	// "BFG9000" => "bfg9000"
	// "BöseÜberraschung" => "böseÜberraschung"
	// "Two  spaces" => "twoSpaces"
	// "BadUTF8\xe2\xe2\xa1" => "badUtf8"
}
