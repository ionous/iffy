/*
 * CODE GENERATED AUTOMATICALLY WITH
 *    github.com/wlbr/templify
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 */

package tables

// runTemplate is a generated function returning the template as a string.
// That string should be parsed by the functions of the golang's template package.
func runTemplate() string {
	var tmpl = "/**\n" +
		" * for saving, restoring a player's game session.\n" +
		" */\n" +
		"create table if not exists \n" +
		"\trun_domain( domain text, active int, primary key( domain )); \n" +
		"\n" +
		"create table if not exists \n" +
		"\trun_pair( noun text, relation text, otherNoun text, unique( noun, relation, otherNoun ) ); \n" +
		""
	return tmpl
}
