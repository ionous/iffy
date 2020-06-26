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
		"\n" +
		"/* the initial values of noun fields: noun, field, value, tier\n" +
		"tier is hierarchy depth, more derived is better */\n" +
		"create view\n" +
		"run_value as \n" +
		"/* future: select *, -1 as tier \n" +
		"\tfrom mdl_run to get existing runtime values */\n" +
		"select noun, field, value, 0 as tier \n" +
		"\tfrom mdl_start ms\n" +
		"union all \n" +
		"select noun, field, value, \n" +
		"\tinstr(mk.kind || \",\" || mk.path || \",\", md.kind || \",\") as tier\n" +
		"\tfrom mdl_noun mn\n" +
		"\tjoin mdl_kind mk\n" +
		"\tusing (kind)\n" +
		"\tjoin mdl_default md\n" +
		"\twhere (tier > 0)\n" +
		"union all \n" +
		"select noun, aspect as field, trait as value, null as tier\n" +
		"\tfrom mdl_noun_traits;\n" +
		""
	return tmpl
}
