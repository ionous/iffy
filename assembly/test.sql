/* resolve default ephemera to strings.
 */
create temp view 
asm_default as
	select p.rowid as idEphDefault, nk.name as kind, nf.name as prop, p.value as value
from eph_default p join eph_named nk
	on (p.idNamedKind = nk.rowid)
left join eph_named nf
	on (p.idNamedProp = nf.rowid);

/* resolve value ephemera to strings.
 */
create temp view
asm_value as
	select pv.rowid as idEphValue, nn.name, np.name as prop, pv.value
from eph_value pv join eph_named nn
	on (pv.idNamedNoun = nn.rowid)
left join eph_named np
	on (pv.idNamedProp = np.rowid);	