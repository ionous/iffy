package assembly

import (
	"database/sql"
	"reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dbutil"
)

func determineInitialFields(m *Modeler, db *sql.DB) (err error) {
	var store valueStore
	var curr, last valueInfo
	if e := dbutil.QueryAll(db,
		`select asm.noun, mf.field, mf.type, asm.value
			from asm_noun as asm
		join mdl_field mf
			on (asm.prop = mf.field)
		where instr((
			select mk.kind || "," || mk.path || ","
			from mdl_kind mk 
			join mdl_noun mn
			using (kind)
			where (mn.noun = asm.noun)
		), mf.kind || ",")
		order by noun, field, type`,
		func() (err error) {
			if nv, e := convertField(curr.fieldType, curr.value); e != nil {
				err = e
			} else if last.target != curr.target || last.field != curr.field {
				store.add(last)
				last, last.value = curr, nv
			} else if !reflect.DeepEqual(last.value, nv) {
				err = errutil.Fmt("conflicting values: %s != %s:%T", last.String(), curr.String(), nv)
			}
			return
		},
		&curr.target, &curr.field, &curr.fieldType,
		&curr.value,
	); e != nil {
		err = e
	} else {
		store.add(last)
		err = store.writeInitialFields(m)
	}
	return
}
