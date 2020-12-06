/*
 * CODE GENERATED AUTOMATICALLY WITH
 *    github.com/wlbr/templify
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 */

package tables

// assemblyTemplate is a generated function returning the template as a string.
// That string should be parsed by the functions of the golang's template package.
func assemblyTemplate() string {
	var tmpl = "\n" +
		"/* resolve (and normalize) eph_kind ephemera to plural strings.\n" +
		" */\n" +
		"create temp view \n" +
		"asm_ancestry as\n" +
		"select ak.name as parent, kn.name as kid \n" +
		"from eph_kind ek\n" +
		"join asm_kind ak\n" +
		"\ton (ak.singularId = ek.idNamedParent)\n" +
		"left join eph_named kn\n" +
		"where kn.rowid = ek.idNamedKind;\n" +
		"\n" +
		"/* resolve test ephemera to strings\n" +
		" */\n" +
		"create temp view\n" +
		"asm_check as select * \n" +
		"from eph_check ek \n" +
		"join eph_named kn\n" +
		"\ton (ek.idNamedTest = kn.rowid);\n" +
		"\n" +
		"create temp view\n" +
		"asm_expect as select * \n" +
		"from eph_expect ex \n" +
		"join eph_named kn\n" +
		"\ton (ex.idNamedTest = kn.rowid);\n" +
		"\n" +
		"\n" +
		"/* resolve default ephemera to strings.\n" +
		" */\n" +
		"create temp view \n" +
		"asm_default as\n" +
		"\tselect p.rowid as idEphDefault, kn.name as kind, nf.name as prop, p.value as value\n" +
		"from eph_default p join eph_named kn\n" +
		"\ton (p.idNamedKind = kn.rowid)\n" +
		"left join eph_named nf\n" +
		"\ton (p.idNamedProp = nf.rowid);\n" +
		"\n" +
		"/**\n" +
		" * using plurals table, convert singular named kinds to plural kinds.\n" +
		" */\n" +
		"create temp view \n" +
		"asm_kind as \n" +
		"\tselect mp.many as name, en.idSource, en.offset, en.rowid as singularId \n" +
		"\tfrom eph_named en \n" +
		"\tjoin mdl_plural mp \n" +
		"\t\ton (mp.one= en.name)\n" +
		"\twhere (en.category = 'singular_kind');\n" +
		"\n" +
		"/* the bits of asm_relation which didnt make it into the mdl_pair table.\n" +
		" */\n" +
		"create temp view \n" +
		"asm_mismatch as\n" +
		"select idEphRel, stem, relation, cardinality, firstNoun, firstKind, secondNoun, secondKind\n" +
		"from asm_relation ar\n" +
		"where max(ar.relation, ar.firstKind, ar.secondKind) is null\n" +
		"or case ar.cardinality\n" +
		"\twhen 'one_one' then\n" +
		"\texists(\n" +
		"\t\tselect 1 \n" +
		"\t\tfrom mdl_pair rel \n" +
		"\t\twhere (ar.relation = rel.relation) \n" +
		"\t\tand ((ar.firstNoun = rel.noun) and (ar.secondNoun != rel.otherNoun)\n" +
		"\t\tor (ar.secondNoun = rel.otherNoun) and (ar.firstNoun != rel.noun))\n" +
		"\t)\n" +
		"\twhen 'one_any' then \n" +
		"\texists(\n" +
		"\t\t/* given otherNoun there is only one valid noun */\n" +
		"\t\tselect 1 \n" +
		"\t\tfrom mdl_pair rel \n" +
		"\t\twhere (ar.relation = rel.relation)\n" +
		"\t\tand (ar.secondNoun = rel.otherNoun) \n" +
		"\t\tand (ar.firstNoun != rel.noun)\n" +
		"\t)\n" +
		"\twhen 'any_one' then \n" +
		"\texists(\n" +
		"\t\t/* given noun there is only one valid otherNoun */\n" +
		"\t\tselect 1 \n" +
		"\t\tfrom mdl_pair rel \n" +
		"\t\twhere (ar.relation = rel.relation)\n" +
		"\t\tand (ar.firstNoun = rel.noun) \n" +
		"\t\tand (ar.secondNoun != rel.otherNoun)\n" +
		"\t)\n" +
		"end;\n" +
		"\n" +
		"/* resolve value ephemera to nouns.\n" +
		"\tmatches nouns by partial name, albeit in a preliminary way.\n" +
		"\tfix? could do a search where asm_noun.noun is NULL to determine missing matches.\n" +
		" */\n" +
		"create temp view\n" +
		"asm_noun as \n" +
		"\tselect *, ( \n" +
		"\t\tselect me.noun \n" +
		"\t\tfrom mdl_name as me\n" +
		"\t\twhere UPPER(asm.name) = UPPER(me.name)\n" +
		"\t \torder by me.rank limit 1 \n" +
		"\t) as noun\n" +
		"from asm_value as asm;\n" +
		"\n" +
		"\n" +
		"/* resolve test ephemera to strings\n" +
		" */\n" +
		"create temp view \n" +
		"asm_pattern as \n" +
		"\tselect pn.name as pattern,\n" +
		"\tkn.name as param, \n" +
		"\ttn.name as type, \n" +
		"\tep.affinity as affinity,\n" +
		"\tep.idProg >=0 as decl, \n" +
		"\tep.rowid as ogid,\n" +
		"\tkn.category as cat,\n" +
		"\tep.idProg\n" +
		"from eph_pattern ep\n" +
		"left join eph_named pn\n" +
		"\ton (ep.idNamedPattern = pn.rowid)\n" +
		"left join eph_named kn\n" +
		"\ton (ep.idNamedParam = kn.rowid)\n" +
		"left join eph_named tn\n" +
		"\ton (ep.idNamedType = tn.rowid);\n" +
		"\n" +
		"/**\n" +
		" * link declared patterns to the successfully modeled types\n" +
		" */\n" +
		"create temp view \n" +
		"asm_pattern_decl as \n" +
		"select pattern, param, type, affinity, ogid,\n" +
		"\t( select mk.kind\n" +
		"\tfrom mdl_kind mk \n" +
		"\tjoin mdl_plural mp\n" +
		"\twhere mp.one = type\n" +
		"\tand mp.many=mk.kind ) as kind, \n" +
		"\tidProg, \n" +
		"\tcat\n" +
		"from asm_pattern \n" +
		"where decl = 1 \n" +
		"group by pattern, param\n" +
		"order by ogid;\n" +
		"\n" +
		"\n" +
		"/* resolve relative ephemera to nouns and relations\n" +
		"use left join(s) to return nulls for missing elements \n" +
		" */\n" +
		"create temp view\n" +
		"asm_relation as\n" +
		"select \t\n" +
		"\tidEphRel, \n" +
		"\tstem, relation, cardinality, \n" +
		"\t\n" +
		"\t/* first contains the kind of the user specified noun;\n" +
		"\t\tswapped contains the kind of the relation \n" +
		"\t*/\n" +
		"\tfirst.noun as firstNoun, \n" +
		"\tcase when instr((\n" +
		"\t \t\t\tselect mk.kind || ',' || mk.path || ','\n" +
		"\t\t\t\tfrom mdl_kind mk\n" +
		"\t\t\t\twhere mk.kind = first.kind\n" +
		"\t\t\t),  swapped.firstKind || ',') \n" +
		"\t\t\tthen first.kind \n" +
		"\tend as firstKind,\n" +
		"\n" +
		"\t/* second contains the kind of the other user specified noun;\n" +
		"\t\tswapped contains the other kind of the relation\n" +
		"\t */\n" +
		"\tsecond.noun as secondNoun,\n" +
		"\tcase when instr((\n" +
		"\t \t\t\tselect mk.kind || ',' || mk.path || ','\n" +
		"\t\t\t\tfrom mdl_kind mk\n" +
		"\t\t\t\twhere mk.kind = second.kind\n" +
		"\t\t\t),  swapped.secondKind || ',') \n" +
		"\t\t\tthen second.kind\n" +
		"\tend as secondKind\n" +
		"from (\n" +
		"\tselect \n" +
		"\t\tidEphRel,stem,relation,cardinality,\n" +
		"\t\tcase swap when 1 then secondNoun else firstNoun end as firstNoun,\n" +
		"\t\tcase swap when 1 then firstNoun else secondNoun end as secondNoun,\n" +
		"\t\tcase swap when 1 then otherKind else kind end as firstKind,\n" +
		"\t\tcase swap when 1 then kind else otherKind end as secondKind\n" +
		"\tfrom (\n" +
		"\t\tselect *, (cardinality = 'one_one') and (ar.firstNoun > ar.secondNoun) as swap\n" +
		"\t\t\tfrom asm_relative ar\n" +
		"\t\t\tleft join asm_verb mv\n" +
		"\t\t\t\tusing (stem)\n" +
		"\t\t\tleft join mdl_rel mr\n" +
		"\t\t\t\tusing (relation)\n" +
		"\t)\n" +
		") as swapped\n" +
		"left join mdl_noun first\n" +
		"\t on (first.noun = swapped.firstNoun)\n" +
		"left join mdl_noun second \n" +
		"\ton (second.noun = swapped.secondNoun);\n" +
		"\n" +
		"/* resolve relative ephemera to nouns.\n" +
		" */\n" +
		"create temp view\n" +
		"asm_relative as \n" +
		"select ar.idEphRel, \n" +
		"\t( select me.noun from mdl_name me\n" +
		"\t\twhere (me.name=ar.firstName) \n" +
		"\t\torder by rank limit 1)\n" +
		"\tas firstNoun, \n" +
		"\tar.stem, \n" +
		"\t( select me.noun from mdl_name me\n" +
		"\t\twhere (me.name=ar.secondName) \n" +
		"\t\torder by rank limit 1)\n" +
		"\tas secondNoun \n" +
		"from asm_relative_name ar\n" +
		"where firstNoun is not null and secondNoun is not null;\n" +
		"\n" +
		"/* resolve relative ephemera to strings.\n" +
		" */\n" +
		"create temp view\n" +
		"asm_relative_name as\n" +
		"select rel.rowid as idEphRel, \n" +
		"\tna.name as firstName, \n" +
		"\tnv.name as stem,\n" +
		"\tnb.name as secondName\n" +
		"from eph_relative rel\n" +
		"join eph_named na\n" +
		"\ton (rel.idNamedHead = na.rowid)\n" +
		"left join eph_named nv\n" +
		"\t\ton (rel.idNamedStem = nv.rowid)\n" +
		"\tleft join eph_named nb\n" +
		"\t\ton (rel.idNamedDependent = nb.rowid);\n" +
		"\n" +
		"\n" +
		"\n" +
		"/* resolve rules to programs\n" +
		" */\n" +
		"create temp view \n" +
		"asm_rule as \n" +
		"\tselect rn.name as pattern, progType as type, prog\n" +
		"from eph_rule er\n" +
		"join eph_named rn\n" +
		"\ton (er.idNamedPattern = rn.rowid)\n" +
		"join eph_prog ep\n" +
		"\ton (er.idProg = ep.rowid)\n" +
		"order by pattern, type, idProg;\n" +
		"\n" +
		"/* patterns and rules with similar names and possibly different types\n" +
		"* fix: does this need to be updated with affinity?\n" +
		" */\n" +
		"create temp view \n" +
		"asm_rule_match as \n" +
		"\tselect pattern, ap.type pt, ar.type rt, prog,\n" +
		"\treplace(ap.type, '_eval', '') =\n" +
		"\treplace(ar.type, '_rule', '') as matched\n" +
		"from asm_rule ar \n" +
		"join asm_pattern ap \n" +
		"using (pattern)\n" +
		"where ap.decl = 1 \n" +
		"and ap.pattern = ap.param\n" +
		"and ap.pattern = ar.pattern;\n" +
		"\n" +
		"/* resolve value ephemera to strings.\n" +
		" */\n" +
		"create temp view\n" +
		"asm_value as\n" +
		"\tselect pv.rowid as idEphValue, nn.name, np.name as prop, pv.value\n" +
		"from eph_value pv join eph_named nn\n" +
		"\ton (pv.idNamedNoun = nn.rowid)\n" +
		"left join eph_named np\n" +
		"\ton (pv.idNamedProp = np.rowid);\n" +
		"\n" +
		"\n" +
		"/* a verb stem implies a specific relation */\n" +
		"create temp table \n" +
		"asm_verb(relation text, stem text, unique(stem));\n" +
		""
	return tmpl
}
