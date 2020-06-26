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
		"create table mdl_aspect(aspect text, trait text, rank int, primary key(aspect, trait));\n" +
		"/* stored tests, a work in progress \n" +
		" * tbd: should named test be turned back into text?\n" +
		" * still need to determine if \"sources\" should be listed in model ( for debugging )\n" +
		" */ \n" +
		"create table mdl_check( name text, idProg int, expect text );\n" +
		"/* default values for the field of a kind ( and its descendant kinds ) */\n" +
		"create table mdl_default(kind text, field text, value blob );\n" +
		"/* properties of a kind. type is a PRIM_ */\n" +
		"create table mdl_field(kind text, field text, type text, primary key(kind, field));\n" +
		"/* a class of objects with shared characteristics */\n" +
		"create table mdl_kind(kind text, path text, primary key(kind));\n" +
		"/* words which refer to nouns. in cases where two words may refer to the same noun, \n" +
		"   the lower rank of the association wins. */\n" +
		"create table mdl_name(noun text, name text, rank int);\n" +
		"/* a person, place, or thing in the world. */\n" +
		"create table mdl_noun(noun text, kind text, primary key(noun));\n" +
		"/* relation between two specific nouns. these change over the course of a game. */\n" +
		"create table mdl_pair(noun text, relation text, otherNoun text);\n" +
		"/* maps common and uncommon words to their plurals */\n" +
		"create table mdl_plural( one text, many text );\n" +
		"/* stored programs, a work in progress; currently copied from eph_prog */ \n" +
		"create table mdl_prog( type text, bytes blob );\n" +
		"/* pattern name and parameter ordering */\n" +
		"create table mdl_pat( pattern text, param text, type text, idx int );\n" +
		"/* relation and constraint between two kinds of nouns */\n" +
		"create table mdl_rel(relation text, kind text, cardinality text, otherKind text, primary key(relation));\n" +
		"/* pattern name and reference to program */\n" +
		"create table mdl_rule( pattern text, idProg int );\n" +
		"/* documentation for pieces of the model: kinds, nouns, fields, etc. */\n" +
		"create table mdl_spec(type text, name text, spec text, primary key(type, name));\n" +
		"/* initial values for various noun properties. these change over the course of a game. */\n" +
		"create table mdl_start(noun text, field text, value blob);\n" +
		"\n" +
		"/**\n" +
		" * all of the traits associated with each of the nouns\n" +
		" * related, fix? should this be \"order by noun, aspect, rank, mt.rowid\"\n" +
		" */\n" +
		"create view \n" +
		"mdl_noun_traits as \n" +
		"select noun, aspect, trait\n" +
		"from mdl_noun \n" +
		"join mdl_kind mk \n" +
		"\tusing (kind)\n" +
		"join mdl_field mf\n" +
		"\ton (mf.type = 'aspect' and \n" +
		"\tinstr((select mk.kind || \",\" || mk.path || \",\"),  mf.kind || \",\"))\n" +
		"join mdl_aspect ma \n" +
		"\ton (ma.aspect = mf.field);"
	return tmpl
}
