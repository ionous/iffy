/*
 * CODE GENERATED AUTOMATICALLY WITH
 *    github.com/wlbr/templify
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 */

package tables

// ephemeraTemplate is a generated function returning the template as a string.
// That string should be parsed by the functions of the golang's template package.
func ephemeraTemplate() string {
	var tmpl = "create table eph_alias( idNamedAlias int, idNamedActual int );\n" +
		"create table eph_aspect( idNamedAspect int );\n" +
		"create table eph_certainty( certainty text, idNamedTrait int, idNamedKind text );\n" +
		"create table eph_check( idNamedTest text, idProg int, expect text );\n" +
		"create table eph_default( idNamedKind int, idNamedProp int, value blob );\n" +
		"create table eph_kind( idNamedKind int, idNamedParent int );\n" +
		"create table eph_named( name text, category text, idSource int, offset text );\n" +
		"create table eph_noun( idNamedNoun int, idNamedKind int );\n" +
		"create table eph_plural( idNamedPlural int, idNamedSingluar int );\n" +
		"create table eph_primitive( primType text, idNamedKind int, idNamedField int );\n" +
		"create table eph_relation( idNamedRelation int, idNamedKind int, idNamedOtherKind int, cardinality text check (cardinality in ('one_one','one_any','any_one','any_any')));\n" +
		"create table eph_relative( idNamedHead int, idNamedStem int, idNamedDependent int );\n" +
		"create table eph_source( src text );\n" +
		"create table eph_trait( idNamedTrait int, idNamedAspect int, rank int );\n" +
		"create table eph_value( idNamedNoun int, idNamedProp int, value blob );\n" +
		"create table eph_verb( idNamedStem int, idNamedRelation int, verb text );\n" +
		"create table eph_prog( idSource int, type text, prog blob );\n" +
		""
	return tmpl
}
