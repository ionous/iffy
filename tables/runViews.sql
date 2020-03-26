/* runtime views */
create temp view
run_init as
select * from (
	/* future: select *, -1 as tier 
		from mdl_run */
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
);