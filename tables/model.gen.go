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
		"/* a class of objects with shared characteristics */\n" +
		"create table mdl_kind(kind text, path text, primary key(kind));\n" +
		"/* properties of a kind */\n" +
		"create table mdl_field(kind text, field text, type text, primary key(kind, field));\n" +
		"/* default values for the field of a kind ( or one of its descendant kinds ) */\n" +
		"create table mdl_default(kind text, field text, value blob );\n" +
		"/* documentation for pieces of the model: kinds, nouns, fields, etc. */\n" +
		"create table mdl_spec(type text, name text, spec string, primary key(type, name));\n" +
		"/* words which refer to nouns. in cases where two words may refer to the same noun, \n" +
		"   the lower rank of the association wins. */\n" +
		"create table mdl_name(noun text, name text, rank int );\n" +
		"/* a person, place, or thing in the world. */\n" +
		"create table mdl_noun(noun text, kind text, primary key(noun));\n" +
		"/* relation and constraint between two kinds of nouns */\n" +
		"create table mdl_rel(relation text, kind text, cardinality text, otherKind text, primary key(relation));\n" +
		"/* relation between two specific nouns */\n" +
		"create table mdl_pair(noun text, relation text, otherNoun text );\n" +
		"/* enumerated value */\n" +
		"create table mdl_aspect(aspect text, trait text, rank int, primary key(aspect, trait));\n" +
		"create table mdl_start(noun text, field text, value blob );\n" +
		""
	return tmpl
}
