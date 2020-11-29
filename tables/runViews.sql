/* runtime views */

/** all of the fields for all of the nouns, including in which kind the field was declared */
create temp view 
mdl_noun_field as 
select noun, field, 
	/* the runtime interprets type as runtime storage type, and aspects are stored as text.
	   fix? change the assembler to output what the runtime expects? */
	case mf.type when 'aspect' then 'text' else mf.type end as type, 
	mk.kind,  
	(mk.kind || ',' || mk.path || ',') as fullpath
from mdl_noun 
join mdl_kind mk 
	using (kind)
join mdl_field mf
	on (instr(fullpath,  mf.kind || ','));

/**
 * all of the traits associated with each of the nouns
 */
create temp view 
mdl_noun_traits as 
select noun, aspect, trait, fullpath
from mdl_noun_field mf
join mdl_aspect ma 
where ma.aspect = mf.field
order by noun, aspect, rank;

/**
 * relational pairs, and the cardinality of their relations
 */
create temp view 
mdl_pair_rel as 
select noun, otherNoun, relation, cardinality, domain
from mdl_pair mp 
join mdl_rel mr 
	using (relation);

/* the initial values of noun fields: noun, field, type, value, tier
tier is hierarchy depth, more derived is better.
more derived classes are on the left, root is on the right, so a small tier is better.
fix? bake down all or part of this table during assembly?
*/
create temp view
run_value as 
/* future: select from mdl_run to get a save game's stored runtime values */
/* for now, the mdl_start fields get the lowest possible tier */
select noun, field, type, value, 0 as tier 
from mdl_noun_field mf
join mdl_start ms 
	using (noun, field)
union all 
select noun, field, type, value, 
	instr(mf.fullpath, md.kind || ',') as tier
	from mdl_noun_field mf
	join mdl_default md
		using (field)
	where (tier > 0) /* greater than zero means this default kind value applies to this noun */	
union all  /* add all default traits at the largest possible tier */
select noun, aspect as field, 'text' as type, trait as value, length(fullpath) as tier
	from mdl_noun_traits
union all  /* add all fields with zero values as null, and we'll order by nulls last */
select noun, field,  type, null as value, null as tier
	from mdl_noun_field;

/* active nouns */
create temp view
run_noun as 
select mn.noun, mn.kind 
from mdl_noun mn 
join run_domain rd 
	on (mn.noun like ('#' || rd.domain || '::%'))
where rd.active = true;