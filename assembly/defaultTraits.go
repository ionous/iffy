package assembly

import (
	"database/sql"
	"fmt"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dbutil"
)

func determineDefaultTraits(m *Modeler, db *sql.DB) (err error) {
	var store traitStore
	var curr, last traitInfo
	if e := dbutil.QueryAll(db,
		// normalize aspect and trait requests
		// we have to do traits and aspects at the same time because
		// they talk about the same pool of values, and could generate conflicts.
		`select asm.kind, mt.aspect, mt.trait,
			ifnull(nullif(asm.value, mt.trait), 1)
		from asm_default as asm
		join mdl_aspect mt
			on (asm.prop = mt.trait) 
			or (asm.prop = mt.aspect and asm.value= mt.trait)
		join mdl_kind mk
			using (kind)
		join mdl_field mf
			on (mf.type = 'aspect')
			and (mf.field = mt.aspect)
		where instr((
			/* path of the noun's kind should contain the kind which declared the aspect*/
			select mk.kind || "," || mk.path || ","
			),  mf.kind || ",")
		order by asm.kind, mt.aspect, mt.trait`,
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
		err = store.writeDefaultTraits(m)
	}
	return
}

type traitInfo struct {
	target, aspect, trait string
	value                 bool
}

func (n *traitInfo) String() string {
	return n.target + "." + n.aspect + ":" + n.trait + fmt.Sprintf("(%v:%T)", n.value, n.value)
}

type traitStore struct {
	list []traitInfo
}

func (store *traitStore) add(n traitInfo) {
	if len(n.target) > 0 {
		store.list = append(store.list, n)
	}
}

func (store *traitStore) writeDefaultTraits(m *Modeler) (err error) {
	for _, n := range store.list {
		if e := m.WriteDefault(n.target, n.aspect, n.trait); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func (store *traitStore) writeInitialTraits(m *Modeler) (err error) {
	for _, n := range store.list {
		if e := m.WriteValue(n.target, n.aspect, n.trait); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}
