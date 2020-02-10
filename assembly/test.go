/*
 * CODE GENERATED AUTOMATICALLY WITH
 *    github.com/wlbr/templify
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 */

package assembly

// testTemplate is a generated function returning the template as a string.
// That string should be parsed by the functions of the golang's template package.
func testTemplate() string {
	var tmpl = "/* resolve default ephemera to strings.\n" +
		" */\n" +
		"create temp view \n" +
		"asm_default as\n" +
		"\tselect p.rowid as idEphDefault, nk.name as kind, nf.name as prop, p.value as value\n" +
		"from eph_default p join eph_named nk\n" +
		"\ton (p.idNamedKind = nk.rowid)\n" +
		"left join eph_named nf\n" +
		"\t\ton (p.idNamedProp = nf.rowid);\n" +
		"\n" +
		"/* resolve value ephemera to strings.\n" +
		" */\n" +
		"create temp view\n" +
		"asm_value as\n" +
		"\tselect pv.rowid as idEphValue, nn.name, np.name as prop, pv.value\n" +
		"from eph_value pv join eph_named nn\n" +
		"\ton (pv.idNamedNoun = nn.rowid)\n" +
		"left join eph_named np\n" +
		"\ton (pv.idNamedProp = np.rowid);\t"
	return tmpl
}
