package test

import "github.com/ionous/iffy/tables"

var FieldData = struct {
	// kind hierarchy
	PathsOfKind,
	// parents of nouns
	KindsOfNoun,
	// noun hierarchy
	PathsOfNoun,
	// kind, field, type
	Fields,
	// aspect, trait pairs
	Traits []string
	// noun, field, value triplets
	DefaultValues, StartingValues,
	// computed noun, field, text value triplets
	TextValues,
	// computed noun, field, num value triplets
	NumValues,
	BoolValues []interface{}
}{
	/*paths of Kind*/ []string{
		"Ks", "",
		"Js", "Ks",
		"Ls", "Ks",
		"Fs", "Ls,Ks",
	},
	/*kinds of noun*/ []string{
		"apple", "Ks",
		"duck", "Js",
		"toy boat", "Ls",
		"boat", "Fs",
	},
	/*paths of noun*/ []string{
		"apple", "Ks",
		"duck", "Js,Ks",
		"toy boat", "Ls,Ks",
		"boat", "Fs,Ls,Ks",
	},
	/*fields*/ []string{
		"Ks", "d", tables.PRIM_DIGI,
		"Ks", "t", tables.PRIM_TEXT,
		"Ks", "a", tables.PRIM_ASPECT,
		"Ls", "b", tables.PRIM_ASPECT,
	},
	/*traits*/ []string{
		"a", "w",
		"a", "x",
		"a", "y",
		"b", "z",
		"b", "zz",
	},
	/*default values*/ []interface{}{
		"Ks", "d", 42,
		"Js", "t", "chippo",
		"Ls", "t", "weazy",
		"Fs", "d", 13,
		"Fs", "b", "zz",
		"Ls", "a", "x",
	},
	/*starting values*/ []interface{}{
		"apple", "d", 5,
		"duck", "d", 1,
		"toy boat", "t", "boboat",
		"boat", "t", "xyzzy",
		"toy boat", "a", "y",
	},
	/*text values*/ []interface{}{
		"apple", "t", "",
		"boat", "t", "xyzzy",
		"duck", "t", "chippo",
		"toy boat", "t", "boboat",
		//
		"apple" /*   */, "a", "w",
		"duck" /*    */, "a", "w",
		"toy boat" /**/, "a", "y",
		"boat" /* */, "a", "x",
		//
		"toy boat" /**/, "b", "z",
		"boat" /* */, "b", "zz",

		// asking for an improper or invalid aspect returns nothing
		// fix? should it return or log error instead?
		"apple" /*   */, "b", nil,
		"boat" /*   */, "G", nil,
	},
	/*num values*/ []interface{}{
		"apple", "d", 5.0,
		"boat", "d", 13.0,
		"duck", "d", 1.0,
		"toy boat", "d", 42.0,
	},
	// noun, truth values. the first comma separated value is true, the rest false.
	/*bool values*/ []interface{}{
		"apple", "w,x,y",
		"duck", "w,x,y",
		//
		"toy boat", "y,w,x",
		"toy boat", "z,zz",
		//
		"boat", "x,w,y",
		"boat", "zz,z",
	},
}
