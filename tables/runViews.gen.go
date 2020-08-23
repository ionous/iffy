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
		"/**\n" +
		" * all of the traits associated with each of the nouns\n" +
		" * note: could use sqlite's special group by behavior to reduce to just the first aspect for each noun.\n" +
		" *  ie. \"group by noun,aspect\"\n" +
		" * -- the run_value clause would work just fine without it.\n" +
		" * related, fix? should this be \"order by noun, aspect, rank, mt.rowid\"\n" +
		" */\n" +
		"create temp view \n" +
		"mdl_noun_traits as \n" +
		"select noun, aspect, trait\n" +
		"from mdl_noun \n" +
		"join mdl_kind mk \n" +
		"\tusing (kind)\n" +
		"join mdl_field mf\n" +
		"\ton (mf.type = 'aspect' and \n" +
		"\tinstr((select mk.kind || \",\" || mk.path || \",\"),  mf.kind || \",\"))\n" +
		"join mdl_aspect ma \n" +
		"\ton (ma.aspect = mf.field)\n" +
		"order by noun, aspect, ma.rank;\n" +
		"\n" +
		"/* the initial values of noun fields: noun, field, value, tier\n" +
		"tier is hierarchy depth, more derived is better */\n" +
		"create temp view\n" +
		"run_value as \n" +
		"/* future: select *, -1 as tier \n" +
		"\tfrom mdl_run to get existing runtime values */\n" +
		"select noun, field, value, 0 as tier \n" +
		"\tfrom mdl_start ms\n" +
		"union all \n" +
		"select noun, field, value, \n" +
		"\tinstr(mk.kind || \",\" || mk.path || \",\", md.kind || \",\") as tier\n" +
		"\tfrom mdl_noun\n" +
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
