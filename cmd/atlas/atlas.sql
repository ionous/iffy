/* any default values defined for the kind;
 null spec indicates the field isnt declared in this kind */
create temp view
atlas_fields as
select kind, field, value, null as spec
	from mdl_default md 
	where not exists (
		select 1 
		from mdl_field mf 
		where mf.kind = md.kind 
		and mf.field = md.field 
	)
union all 
/* and all of the fields defined for the kind */
select 
	kind, 
	field, 
	coalesce((
	/* with the default specified value */
		select value 
		from mdl_default md 
		where mf.kind = md.kind 
		and mf.field = md.field 
		limit 1
		),
	/* or, use type-dependent default value */
	case mf.type 
		when 'aspect' then (
			select trait 
			from mdl_aspect 
			where aspect = field
			order by rank desc
			limit 1
		)
		when 'digi' then '0'
		when 'text' then '""'
		else '???'||mf.type
	end)
	as value, 
	/* include the spec */
	coalesce((
		select spec from mdl_spec spec
		where (spec.type = 'field'
		and spec.name = (kind||'.'||field))
		limit 1 ), '')
	as spec
from mdl_field mf;
	