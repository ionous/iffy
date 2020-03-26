/*
 * CODE GENERATED AUTOMATICALLY WITH
 *    github.com/wlbr/templify
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 */

package tables

// runViewsTemplate is a generated function returning the template as a string.
// That string should be parsed by the functions of the golang's template package.
func runViewsTemplate() string {
	var tmpl = "/* runtime views */\n" +
		"create temp view\n" +
		"run_init as\n" +
		"select * from (\n" +
		"\t/* future: select *, -1 as tier \n" +
		"\t\tfrom mdl_run */\n" +
		"\tselect noun, field, value, 0 as tier\n" +
		"\t\tfrom mdl_start ms\n" +
		"\tunion all \n" +
		"\tselect noun, field, value, \n" +
		"\t\tinstr(mk.kind || \",\" || mk.path || \",\", md.kind || \",\") as tier\n" +
		"\t\tfrom mdl_noun mn\n" +
		"\t\tjoin mdl_kind mk\n" +
		"\t\tusing (kind)\n" +
		"\t\tjoin mdl_default md\n" +
		"\t\twhere (tier > 0)\n" +
		");"
	return tmpl
}
