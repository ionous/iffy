package assembly

import "database/sql"

// goal: build tables of noun, relation, otherNoun
func DetermineRelatives(db *sql.DB) error {
	_, e := db.Exec(`
insert into mdl_pair( noun, relation, otherNoun )
select distinct noun, relation, otherNoun from (
	select *, row_number() over n1 as n1, row_number() over n2 as n2
	from asm_relation
	where max(noun, stem, otherNoun, relation, kind, otherKind) is not null
	window 
	/* count the times the nouns appear in their respective columns */
      n1 as (partition by relation,noun),
	  n2 as (partition by relation,otherNoun)
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
