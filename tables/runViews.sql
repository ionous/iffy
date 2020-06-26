/* runtime views */

/* the initial values of noun fields: noun, field, value, tier
tier is hierarchy depth, more derived is better */
create view
run_value as 
/* future: select *, -1 as tier 
	from mdl_run to get existing runtime values */
select noun, field, value, 0 as tier 
	from mdl_start ms
union all 
select noun, field, value, 
	instr(mk.kind || "," || mk.path || ",", md.kind || ",") as tier
	from mdl_noun mn
	join mdl_kind mk
	using (kind)
	join mdl_default md
	where (tier > 0)
union all 
select noun, aspect as field, trait as value, null as tier
	from mdl_noun_traits;
