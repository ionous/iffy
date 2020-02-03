package assembly

import "database/sql"

// goal: build tables of noun, relation, otherNoun
func DetermineRelatives(db *sql.DB) error {
	_, e := db.Exec(`
insert into start_rel( noun, relation, otherNoun )
select distinct noun, relation, otherNoun from (
	select *, rank() over n1 as n1, rank() over n2 as n2
	from asm_relation a 
	where (a.cardinality != 'one_one')
	or (a.noun = a.otherNoun)
	/* exclude cross column results for 1-1 relations */
	or not exists (
		select 1 from asm_relation b 
		where (a.relation = b.relation)
		and (a.noun = b.otherNoun)
	)
	window 
	/* count the number of times the noun, other noun appears */
      n1 AS (partition by noun order by idEphRel),
	  n2 AS (partition by otherNoun order by idEphRel)
) 
where case cardinality
	when 'one_one'
		then n1 = 1 and n2 = 1
	when 'one_any'
		then n2 = 1
	when 'any_one'
		then n1 = 1
	else 1
end`)
	return e
}
