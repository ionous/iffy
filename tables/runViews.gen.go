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
		"/** all of the fields for all of the nouns, including in which kind the field was declared */\n" +
		"create temp view \n" +
		"mdl_noun_field as \n" +
		"select noun, field, \n" +
		"\t/* the runtime interprets type as runtime storage type, and aspects are stored as text.\n" +
		"\t   fix? change the assembler to output what the runtime expects? */\n" +
		"\tcase mf.type when 'aspect' then 'text' else mf.type end as type, \n" +
		"\tmk.kind,  \n" +
		"\t(mk.kind || ',' || mk.path || ',') as fullpath\n" +
		"from mdl_noun \n" +
		"join mdl_kind mk \n" +
		"\tusing (kind)\n" +
		"join mdl_field mf\n" +
		"\ton (instr(fullpath,  mf.kind || ','));\n" +
		"\n" +
		"/**\n" +
		" * all of the traits associated with each of the nouns\n" +
		" */\n" +
		"create temp view \n" +
		"mdl_noun_traits as \n" +
		"select noun, aspect, trait, fullpath\n" +
		"from mdl_noun_field mf\n" +
		"join mdl_aspect ma \n" +
		"where ma.aspect = mf.field\n" +
		"order by noun, aspect, rank;\n" +
		"\n" +
		"/**\n" +
		" * relational pairs, and the cardinality of their relations\n" +
		" */\n" +
		"create temp view \n" +
		"mdl_pair_rel as \n" +
		"select noun, otherNoun, relation, cardinality, domain\n" +
		"from mdl_pair mp \n" +
		"join mdl_rel mr \n" +
		"\tusing (relation);\n" +
		"\n" +
		"/* the initial values of noun fields: noun, field, type, value, tier\n" +
		"tier is hierarchy depth, more derived is better.\n" +
		"more derived classes are on the left, root is on the right, so a small tier is better.\n" +
		"fix? bake down all or part of this table during assembly?\n" +
		"*/\n" +
		"create temp view\n" +
		"run_value as \n" +
		"/* future: select from mdl_run to get a save game's stored runtime values */\n" +
		"/* for now, the mdl_start fields get the lowest possible tier */\n" +
		"select noun, field, type, value, 0 as tier \n" +
		"from mdl_noun_field mf\n" +
		"join mdl_start ms \n" +
		"\tusing (noun, field)\n" +
		"union all \n" +
		"select noun, field, type, value, \n" +
		"\tinstr(mf.fullpath, md.kind || ',') as tier\n" +
		"\tfrom mdl_noun_field mf\n" +
		"\tjoin mdl_default md\n" +
		"\t\tusing (field)\n" +
		"\twhere (tier > 0) /* greater than zero means this default kind value applies to this noun */\t\n" +
		"union all  /* add all default traits at the largest possible tier */\n" +
		"select noun, aspect as field, 'text' as type, trait as value, length(fullpath) as tier\n" +
		"\tfrom mdl_noun_traits\n" +
		"union all  /* add all fields with zero values as null, and we'll order by nulls last */\n" +
		"select noun, field,  type, null as value, null as tier\n" +
		"\tfrom mdl_noun_field;\n" +
		""
	return tmpl
}
