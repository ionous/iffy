/*
 * CODE GENERATED AUTOMATICALLY WITH
 *    github.com/wlbr/templify
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 */

package tables

// modelTemplate is a generated function returning the template as a string.
// That string should be parsed by the functions of the golang's template package.
func modelTemplate() string {
	var tmpl = "/**\n" +
		" * tables describing the the world and its rules.\n" +
		" */\n" +
		"\n" +
		"/* enumerated values used by kinds and nouns */\n" +
		"create table mdl_aspect( aspect text, trait text, rank int, primary key( aspect, trait ));\n" +
		"/* stored tests, a work in progress \n" +
		" * ex. determine if \"sources\" should be listed in model ( for debugging )\n" +
		" */ \n" +
		"create table mdl_check( name text, type text, expect text );\n" +
		"/* default values for the field of a kind ( and its descendant kinds ) */\n" +
		"create table mdl_default( kind text, field text, value blob );\n" +
		"/* todo: hierarchy of domains \n" +
		"create table mdl_domain( domain text, path text, primary key( domain ));\n" +
		"*/\n" +
		"/* properties of a kind. type is a PRIM_ */\n" +
		"create table mdl_field( kind text, field text, type text, primary key( kind, field ));\n" +
		"/* a class of objects with shared characteristics */\n" +
		"create table mdl_kind( kind text, path text, primary key( kind ));\n" +
		"/* words which refer to nouns. in cases where two words may refer to the same noun, \n" +
		"   the lower rank of the association wins. */\n" +
		"create table mdl_name( noun text, name text, rank int );\n" +
		"/* a person, place, or thing in the world. */\n" +
		"create table mdl_noun( noun text, kind text, primary key( noun ));\n" +
		"/* relation between two specific nouns. these change over the course of a game. */\n" +
		"create table mdl_pair( noun text, relation text, otherNoun text );\n" +
		"/* maps common and uncommon words to their plurals */\n" +
		"create table mdl_plural( one text, many text );\n" +
		"/* stored programs, a work in progress \n" +
		" * note: name is not necessarily unique. some programs are stored in parts:\n" +
		"   for example, patterns and tests. tbd whether the tables should self-describe,\n" +
		"   the connection between tests and patterns and progs are, essentially, application knowledge. */ \n" +
		"create table mdl_prog( name text, type text, bytes blob );\n" +
		"/* declared patterns and their parameters -- \n" +
		"   used mainly for translating pattern positional parameters into names */\n" +
		"create table mdl_pat( pattern text, param text, type text, idx int );\n" +
		"/* relation and constraint between two kinds of nouns */\n" +
		"create table mdl_rel( relation text, kind text, cardinality text, otherKind text, primary key( relation ));\n" +
		"/* documentation for pieces of the model: kinds, nouns, fields, etc. */\n" +
		"create table mdl_spec( type text, name text, spec text, primary key( type, name ));\n" +
		"/* initial values for various noun properties. these change over the course of a game. */\n" +
		"create table mdl_start( noun text, field text, value blob );"
	return tmpl
}
