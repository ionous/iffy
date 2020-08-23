/* runtime views */
/**
 * all of the traits associated with each of the nouns
 * note: could use sqlite's special group by behavior to reduce to just the first aspect for each noun.
 *  ie. "group by noun,aspect"
 * -- the run_value clause would work just fine without it.
 * related, fix? should this be "order by noun, aspect, rank, mt.rowid"
 */
create temp view 
mdl_noun_traits as 
select noun, aspect, trait
from mdl_noun 
join mdl_kind mk 
	using (kind)
join mdl_field mf
	on (mf.type = 'aspect' and 
	instr((select mk.kind || "," || mk.path || ","),  mf.kind || ","))
join mdl_aspect ma 
	on (ma.aspect = mf.field)
order by noun, aspect, ma.rank;

/* the initial values of noun fields: noun, field, value, tier
tier is hierarchy depth, more derived is better */
create temp view
run_value as 
/* future: select *, -1 as tier 
	from mdl_run to get existing runtime values */
select noun, field, value, 0 as tier 
	from mdl_start ms
union all 
select noun, field, value, 
	instr(mk.kind || "," || mk.path || ",", md.kind || ",") as tier
	from mdl_noun
	join mdl_kind mk
	using (kind)
	join mdl_default md
	where (tier > 0)
union all 
select noun, aspect as field, trait as value, null as tier
	from mdl_noun_traits;
