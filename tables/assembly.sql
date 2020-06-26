/* resolve test ephemera to strings
 */
create temp view 
asm_pattern as 
	select pn.name as pattern, 
	kn.name as param, 
	tn.name as type, 
	ep.decl, 
	ep.rowid as ogid
from eph_pattern ep
left join eph_named pn
	on (ep.idNamedPattern = pn.rowid)
left join eph_named kn
	on (ep.idNamedParam = kn.rowid)
left join eph_named tn
	on (ep.idNamedType = tn.rowid);

create temp view 
asm_pattern_decl as 
	select pattern, param, type, ogid 
	from asm_pattern 
	where decl = 1 
	group by pattern, param
	order by ogid;

/**
 * using plurals table, convert singular named kinds to plural kinds.
 */
create temp view 
asm_kind as 
	select mp.many as name, en.idSource, en.offset, en.rowid as singularId 
	from eph_named en 
	join mdl_plural mp 
		on (mp.one= en.name)
	where (en.category = 'singular_kind');

/* resolve (and normalize) eph_kind ephemera to plural strings.
 */
create temp view 
asm_ancestry as
select ak.name as parent, nk.name as kid 
from eph_kind ek
join asm_kind ak
	on (ak.singularId = ek.idNamedParent)
left join eph_named nk
where nk.rowid = ek.idNamedKind;

/* resolve test ephemera to strings
 */
create temp view 
asm_rule as 
	select rn.name as pattern, type, prog, idProg
from eph_rule er
left join eph_named rn
	on (er.idNamedPattern = rn.rowid)
left join eph_prog ep
	on (er.idProg = ep.rowid);

/* patterns and rules with similar names and possibly different types
 */
create temp view 
asm_rule_match as 
	select pattern, ap.type pt, ar.type rt, prog,
	replace(ap.type, '_eval', '') =
	replace(ar.type, '_rule', '') as matched
from asm_rule ar 
join asm_pattern ap 
using (pattern)
where ap.decl = 1 
and ap.pattern = ap.param
and ap.pattern = ar.pattern;

/* resolve test ephemera to strings
 */
create temp view
asm_check as
	select nk.name as name, idProg, expect 
from eph_check p join eph_named nk
	on (p.idNamedTest = nk.rowid);

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

/* resolve value ephemera to nouns.
 */
create temp view
asm_noun as 
	select *, ( 
		select n.noun 
		from mdl_name as n
		where asm.name = n.name 
	 	order by rank
		limit 1 
	) as noun
from asm_value as asm;

/* resolve relative ephemera to strings.
 */
create temp view
asm_relative as
select rel.rowid as idEphRel, 
	na.name as noun, 
	nv.name as stem,
	nb.name as otherNoun
from eph_relative rel
join eph_named na
	on (rel.idNamedHead = na.rowid)
left join eph_named nv
		on (rel.idNamedStem = nv.rowid)
	left join eph_named nb
		on (rel.idNamedDependent = nb.rowid);

/* resolve relative ephemera to nouns and relations
use left join(s) to return nulls for missing elements 
 */
create temp view
asm_relation as
select 	
	idEphRel, 
	stem, relation, cardinality, 
	
	/* first contains the kind of the user specified noun;
		swapped contains the kind of the relation 
	*/
	first.noun as noun, 
	case when instr((
	 			select mk.kind || "," || mk.path || ","
				from mdl_kind mk
				where mk.kind = first.kind
			),  swapped.kind || ",") 
			then first.kind 
	end as kind,

	/* second contains the kind of the other user specified noun;
		swapped contains the other kind of the relation
	 */
	second.noun as otherNoun,
	case when instr((
	 			select mk.kind || "," || mk.path || ","
				from mdl_kind mk
				where mk.kind = second.kind
			),  swapped.otherKind || ",") 
			then second.kind
	end as otherKind
from (
	select 
		idEphRel,stem,relation,cardinality,
		case swap when 1 then otherNoun else noun end as noun,
		case swap when 1 then noun else otherNoun end as otherNoun,
		case swap when 1 then otherKind else kind end as kind,
		case swap when 1 then kind else otherKind end as otherKind
	from (
		select *, (cardinality = 'one_one') and (noun > otherNoun) as swap
			from asm_relative ar
			left join asm_verb mv
				using (stem)
			left join mdl_rel mr
				using (relation)
	)
) as swapped
left join mdl_noun first
	 on (first.noun = swapped.noun)
left join mdl_noun second 
	on (second.noun = swapped.otherNoun);

/* the bits of asm_relation which didnt make it into the mdl_pair table.
 */
create temp view 
asm_mismatch as
select idEphRel, stem, relation, cardinality, noun, kind, otherNoun, otherKind
from asm_relation asm
where max(asm.relation, asm.kind, asm.otherKind) is null
or case asm.cardinality
	when 'one_one' then
	exists(
		select 1 
		from mdl_pair rel 
		where (asm.relation = rel.relation) 
		and ((asm.noun = rel.noun) and (asm.otherNoun != rel.otherNoun)
		or (asm.otherNoun = rel.otherNoun) and (asm.noun != rel.noun))
	)
	when 'one_any' then 
	exists(
		/* given otherNoun there is only one valid noun */
		select 1 
		from mdl_pair rel 
		where (asm.relation = rel.relation)
		and (asm.otherNoun = rel.otherNoun) 
		and (asm.noun != rel.noun)
	)
	when 'any_one' then 
	exists(
		/* given noun there is only one valid otherNoun */
		select 1 
		from mdl_pair rel 
		where (asm.relation = rel.relation)
		and (asm.noun = rel.noun) 
		and (asm.otherNoun != rel.otherNoun)
	)
end;

/* a verb stem implies a specific relation */
create temp table 
asm_verb(relation text, stem text, unique(stem));
