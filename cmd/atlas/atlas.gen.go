/*
 * CODE GENERATED AUTOMATICALLY WITH
 *    github.com/wlbr/templify
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 */

package main

// atlasTemplate is a generated function returning the template as a string.
// That string should be parsed by the functions of the golang's template package.
func atlasTemplate() string {
	var tmpl = "/* any default values defined for the kind;\n" +
		" null spec indicates the field isnt declared in this kind */\n" +
		"create view\n" +
		"atlas_fields as\n" +
		"select kind, field, value, null as spec\n" +
		"\tfrom mdl_default md \n" +
		"\twhere not exists (\n" +
		"\t\tselect 1 \n" +
		"\t\tfrom mdl_field mf \n" +
		"\t\twhere mf.kind = md.kind \n" +
		"\t\tand mf.field = md.field \n" +
		"\t)\n" +
		"union all \n" +
		"/* and all of the fields defined for the kind */\n" +
		"select \n" +
		"\tkind, \n" +
		"\tfield, \n" +
		"\tcoalesce((\n" +
		"\t/* with the default specified value */\n" +
		"\t\tselect value \n" +
		"\t\tfrom mdl_default md \n" +
		"\t\twhere mf.kind = md.kind \n" +
		"\t\tand mf.field = md.field \n" +
		"\t\tlimit 1\n" +
		"\t\t),\n" +
		"\t/* or, use type-dependent default value */\n" +
		"\tcase mf.type \n" +
		"\t\twhen 'aspect' then (\n" +
		"\t\t\tselect trait \n" +
		"\t\t\tfrom mdl_aspect \n" +
		"\t\t\twhere aspect = field\n" +
		"\t\t\torder by rank desc\n" +
		"\t\t\tlimit 1\n" +
		"\t\t)\n" +
		"\t\twhen 'number' then '0'\n" +
		"\t\twhen 'text' then '\"\"'\n" +
		"\t\telse '???'||mf.type\n" +
		"\tend)\n" +
		"\tas value, \n" +
		"\t/* include the spec */\n" +
		"\tcoalesce((\n" +
		"\t\tselect spec from mdl_spec spec\n" +
		"\t\twhere (spec.type = 'field'\n" +
		"\t\tand spec.name = (kind||'.'||field))\n" +
		"\t\tlimit 1 ), '')\n" +
		"\tas spec\n" +
		"from mdl_field mf;\n" +
		"\t"
	return tmpl
}
