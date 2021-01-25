package assembly

// goal: build tables of noun, relation, otherNoun
func AssembleRelatives(asm *Assembler) error {
	_, e := asm.cache.DB().Exec(`
insert into mdl_pair( noun, relation, otherNoun, domain )
select distinct firstNoun, relation, secondNoun, domain from (
	select *, row_number() over n1 as n1, row_number() over n2 as n2
	from asm_relation
	where max(firstNoun, stem, secondNoun, relation, firstKind, secondKind) is not null
	window 
	/* count the times the nouns appear in their respective columns */
      n1 as (partition by relation,firstNoun),
	  n2 as (partition by relation,secondNoun)
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
