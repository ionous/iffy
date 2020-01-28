package assembly

import (
	"database/sql"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dbutil"
	_ "github.com/mattn/go-sqlite3"
)

func determineDefaultTraits(m *Modeler, db *sql.DB) (err error) {
	var store defaultValueStore
	var curr, last defaultValue
	if e := dbutil.QueryAll(db,
		// normalize aspect and trait requests
		`with aspect as (
		select asm.kind, mt.aspect, mt.trait, asm.value 
			from asm_default as asm 
			join mdl_trait mt
				on (asm.prop=mt.trait)
			join mdl_aspect ma
				using (aspect)
		union all 
		select asm.kind, mt.aspect, mt.trait, true as value
			from asm_default as asm 
			join mdl_trait mt
				on (asm.prop=mt.aspect)
				and (asm.value=mt.trait)
		)
		select at.kind, at.aspect, at.trait, at.value from aspect at
		/* filter if the same named trait appears in different aspects */
			where instr((
				select mk.kind || "," || mk.path || ","
				from mdl_kind mk 
				where mk.kind = at.kind
			),  at.kind || ",")
			order by at.kind, at.aspect, at.trait, at.value`,
		func() (err error) {
			if v, ok := curr.value.(int64); !ok || v == 0 {
				// future: re: certainty values
				err = errutil.Fmt("only positive traits are accepted right now")
			} else if last.target != curr.target || last.fieldType != curr.fieldType {
				store.add(last)
				last = curr
			} else if last.field != curr.field {
				err = errutil.Fmt("conflicting defaults: %s != %s", last.String(), curr.String())
			}
			return
		},
		// kind, aspect, trait
		&curr.target, &curr.fieldType, &curr.field,
		&curr.value,
	); e != nil {
		err = e
	} else {
		store.add(last)
		err = store.writeTraits(m)
	}
	return
}

func (store *defaultValueStore) writeTraits(m *Modeler) (err error) {
	for _, n := range store.list {
		if e := m.WriteDefault(n.target, n.fieldType, n.field); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}
