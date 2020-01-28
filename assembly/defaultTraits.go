package assembly

import (
	"database/sql"
	"fmt"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dbutil"
	_ "github.com/mattn/go-sqlite3"
)

func determineDefaultTraits(m *Modeler, db *sql.DB) (err error) {
	var store defaultTraitStore
	var curr, last defaultTrait
	if e := dbutil.QueryAll(db,
		// normalize aspect and trait requests
		`with aspect as (
		select asm.kind, mt.aspect, mt.trait, asm.value 
			from asm_default as asm 
			join mdl_trait mt
				on (asm.prop=mt.trait)
			join mdl_field mf
				on (mf.type = 'aspect')
				and (mf.field = mt.aspect)
		union all 
		select asm.kind, mt.aspect, mt.trait, true as value
			from asm_default as asm 
			join mdl_trait mt
				on (asm.prop = mt.aspect)
				and (asm.value = mt.trait)
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
			if !curr.value {
				// future: possibly a switch for false values that tries to select a single opposite?
				// possibly a separate table for opposites? ( re: relations )
				err = errutil.Fmt("only positive traits are accepted right now")
			} else if last.kind != curr.kind || last.aspect != curr.aspect {
				store.add(last)
				last = curr
			} else if last.trait != curr.trait {
				err = errutil.Fmt("conflicting defaults: %s != %s", last.String(), curr.String())
			}
			return
		},
		&curr.kind, &curr.aspect, &curr.trait,
		&curr.value,
	); e != nil {
		err = e
	} else {
		store.add(last)
		err = store.writeTraits(m)
	}
	return
}

type defaultTrait struct {
	kind, aspect, trait string
	value               bool
}

func (n *defaultTrait) String() string {
	return n.kind + "." + n.aspect + ":" + n.trait + fmt.Sprintf("(%v:%T)", n.value, n.value)
}

type defaultTraitStore struct {
	list []defaultTrait
}

func (store *defaultTraitStore) add(n defaultTrait) {
	if len(n.kind) > 0 {
		store.list = append(store.list, n)
	}
}

func (store *defaultTraitStore) writeTraits(m *Modeler) (err error) {
	for _, n := range store.list {
		if e := m.WriteDefault(n.kind, n.aspect, n.trait); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}
