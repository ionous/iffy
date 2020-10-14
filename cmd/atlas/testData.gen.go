/*
 * CODE GENERATED AUTOMATICALLY WITH
 *    github.com/wlbr/templify
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 */

package main

// testDataTemplate is a generated function returning the template as a string.
// That string should be parsed by the functions of the golang's template package.
func testDataTemplate() string {
	var tmpl = "/* ---------------------------------------------------\n" +
		"test data for web\n" +
		" --------------------------------------------------- */\n" +
		"insert into mdl_kind(kind, path) values\n" +
		"\t('things', ''),\n" +
		"\t('vehicles', 'things'),\n" +
		"\t('cars', 'vehicles,things'),\n" +
		"\t('people', 'things');\n" +
		"\n" +
		"insert into mdl_field(kind, field, type) values\n" +
		"\t('things', 'brief', 'text'),\n" +
		"\t('vehicles', 'flightiness', 'aspect'),\n" +
		"\t('cars', 'num wheels', 'number');\n" +
		"\n" +
		"insert into mdl_default(kind, field, value) values\n" +
		"\t('cars', 'flightiness', 'flightless'),\n" +
		"\t('cars', 'num wheels', 4);\n" +
		"\n" +
		"insert into mdl_aspect(aspect, trait) values\n" +
		"\t('flightiness', 'flightless'),\n" +
		"\t('flightiness', 'glide worthy'),\n" +
		"\t('flightiness', 'flight worthy');\n" +
		"\n" +
		"insert into mdl_rel(relation, kind, cardinality, otherKind) values\n" +
		"\t('containing', 'vehicles', 'one_any', 'people');\n" +
		"\n" +
		"insert into mdl_noun(noun, kind) values\n" +
		"\t('dune buggy', 'cars'),\n" +
		"\t('enterprise', 'vehicles'),\n" +
		"\t('riker', 'people'),\n" +
		"\t('picard', 'people');\n" +
		"\n" +
		"insert into mdl_name(noun, name, rank) values\n" +
		"\t('dune buggy', 'dune buggy', 0),\n" +
		"\t('dune buggy', 'dune', 1),\n" +
		"\t('dune buggy', 'buggy', 1),\n" +
		"\t('enterprise', 'enterprise', 0),\n" +
		"\t('riker', 'riker', 0),\n" +
		"\t('picard', 'picard', 0);\n" +
		"\n" +
		"insert into mdl_start(noun, field, value) values\n" +
		"\t( 'dune buggy', 'num wheels', 3);\n" +
		"\n" +
		"insert into mdl_pair(noun, relation, otherNoun) values\n" +
		"\t( 'dune buggy', 'containing', 'picard'),\n" +
		"\t( 'enterprise', 'containing', 'riker');\n" +
		"\n" +
		"insert into mdl_spec(type, name, spec) values\n" +
		"\t('kind',  'things', 'From inform: ''Represents anything interactive in the world. People, pieces of scenery, furniture, doors and mislaid umbrellas might all be examples, and so might more surprising things like the sound of birdsong or a shaft of sunlight.'''),\n" +
		"\t('aspect', 'flightiness', 'The flight worthiness of vehicles, an example of an aspect with several traits.'),\n" +
		"\t('trait', 'glide worthy', 'Better at landing than taking off.'),\n" +
		"\t('relation', 'containing', 'The outside of insides.'),\n" +
		"\t('field', 'cars.num wheels', 'Not all cars are created equal, or even even.');\n" +
		"\n" +
		""
	return tmpl
}
