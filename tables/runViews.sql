/* runtime views */
create temp view
run_init as
select * from (
	/* future: select *, -1 as tier 
		from mdl_run to get existing runtime values */
	select noun, field, value, 0 as tier /* tier is hierarchy depth, more derived is better */
		from mdl_start ms
	union all 
	select noun, field, value, 
		instr(mk.kind || "," || mk.path || ",", md.kind || ",") as tier
		from mdl_noun mn
		join mdl_kind mk
		using (kind)
		join mdl_default md
		where (tier > 0)
);