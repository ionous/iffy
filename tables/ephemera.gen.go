/*
 * CODE GENERATED AUTOMATICALLY WITH
 *    github.com/wlbr/templify
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 */

package tables

// ephemeraTemplate is a generated function returning the template as a string.
// That string should be parsed by the functions of the golang's template package.
func ephemeraTemplate() string {
	var tmpl = "/* alternative name for a noun */\n" +
		"create table eph_alias( idNamedAlias int, idNamedActual int );\n" +
		"/* collection of related states ( called traits ) */\n" +
		"create table eph_aspect( idNamedAspect int );\n" +
		"/* likelihood that a trait applies to a particular kind */\n" +
		"create table eph_certainty( certainty text, idNamedTrait int, idNamedKind text );\n" +
		"/* test programs and the results they are expected to produce */\n" +
		"create table eph_check( idNamedTest int, idProg int );\n" +
		"/* initial values for the properties of nouns of belonging to the specified kind */\n" +
		"create table eph_default( idNamedKind int, idNamedProp int, value blob );\n" +
		"/* test programs and the results they are expected to produce */\n" +
		"create table eph_expect( idNamedTest int, testType text, expect text );\n" +
		"/* property name and type associated with a kind of object */\n" +
		"create table eph_field( primType text, idNamedKind int, idNamedField int );\n" +
		"/* collection of related nouns, plural named kind, singular named parent\n" +
		" * ex. cats are a kind of animal.\n" +
		" */\n" +
		"create table eph_kind( idNamedKind int, idNamedParent int );\n" +
		"/* user specified appellation and the location that specification came from \n" +
		" * domain references another name\n" +
		" */\n" +
		"create table eph_named( name text, og text, category text, domain int, idSource int, offset text );\n" +
		"/* a named object in the game world and its kind (singular) */\n" +
		"create table eph_noun( idNamedNoun int, idNamedKind int );\n" +
		"/* declarations and references to pattern parameter and pattern return types.\n" +
		"if idNamedPattern is idNamedParam it indicates a return type.\n" +
		"idNamedParam can be in the kind category, or one of the predefined eval types .\n" +
		"idProg holds program initialization for locals ( and possibly default parameter values )\n" +
		"for now idProg < 0 is a pattern or pattern parameter reference;\n" +
		"the \"category\" of idNamedParam also indicates reference vs. declaration\n" +
		"*/\n" +
		"create table eph_pattern( idNamedPattern int, idNamedParam int, idNamedType int, idProg int );\n" +
		"/* rule for the collective name of a singular word */\n" +
		"create table eph_plural( idNamedPlural int, idNamedSingluar int );\n" +
		"/* type is the name of the command container for de-serialization of the prog\n" +
		" * idSource is the origin of the ephemera just like eph_named.\n" +
		" * fix? unclear why there's no offset, and if here is domain, source, offset concerns\n" +
		" * why that isnt a table of its own for declarations to share.\n" +
		" */\n" +
		"create table eph_prog( idSource int, progType text, prog blob );\n" +
		"/* connection between two kinds of object */\n" +
		"create table eph_relation( idNamedRelation int, idNamedKind int, idNamedOtherKind int, cardinality text check (cardinality in ('one_one','one_any','any_one','any_any')));\n" +
		"/* connection between two object instances */\t\n" +
		"create table eph_relative( idNamedHead int, idNamedStem int, idNamedDependent int );\n" +
		"/* function handler for a pattern */\n" +
		"create table eph_rule( idNamedPattern int, idProg int );/* uri, file name or other identification for the origin of the various ephemera recorded in the db. \n" +
		"while its not particularly useful to have a one column primitive data type column\n" +
		"someday, this might contain source modification times or other useful info.\n" +
		"*/\n" +
		"create table eph_source( src text );\n" +
		"/* only one trait from a given aspect can be true for a noun at a time. */\t\n" +
		"create table eph_trait( idNamedTrait int, idNamedAspect int, rank int );\n" +
		"/* initial value for a noun's field, trait, or aspect \n" +
		" * fix? why does there need to be both value and default?\n" +
		" * the names ( idNamedNoun vs idNamedKind ) should be separation enough.\n" +
		" */\n" +
		"create table eph_value( idNamedNoun int, idNamedProp int, value blob );\n" +
		"/* word indicating a particular relationship between nouns */\n" +
		"create table eph_verb( idNamedStem int, idNamedRelation int, verb text );\n" +
		""
	return tmpl
}
