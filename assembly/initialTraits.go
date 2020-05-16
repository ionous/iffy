package assembly

import (
	"database/sql"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/tables"
)

// reads mdl_aspect, mdl_noun, mdl_kind, mdl_field
func determineInitialTraits(m *Modeler, db *sql.DB) (err error) {
	var store traitStore
	var curr, last traitInfo
	if e := tables.QueryAll(db,
		// normalize aspect and trait requests
		`select asm.noun, mt.aspect, mt.trait, 
		 	ifnull(nullif(asm.value, mt.trait), 1)
		from asm_noun as asm 
		join mdl_aspect mt
			on (asm.prop = mt.trait) 
			or (asm.prop = mt.aspect and asm.value= mt.trait)
		join mdl_noun mn
			using (noun)
		join mdl_kind mk
			using (kind)
		join mdl_field mf
			on (mf.type = 'aspect')
			and (mf.field = mt.aspect)
		where instr((
			/* path of the noun's kind should contain the kind which declared the aspect*/
				select mk.kind || "," || mk.path || ","
			), mf.kind || ",")
		order by noun, aspect, trait`,
		func() (err error) {
			if !curr.value {
				// future: possibly a switch for false values that tries to select a single opposite?
				// possibly a separate table for opposites? ( re: relations )
				err = errutil.Fmt("only positive traits are accepted right now")
			} else if last.target != curr.target || last.aspect != curr.aspect {
				store.add(last)
				last = curr
			} else if last.trait != curr.trait {
				err = errutil.Fmt("conflicting defaults: %s != %s", last.String(), curr.String())
			}
			return
		},
		&curr.target, &curr.aspect, &curr.trait,
		&curr.value,
	); e != nil {
		err = e
	} else {
		store.add(last)
		err = store.writeInitialTraits(m)
	}
	return
}
