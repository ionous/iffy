package assembly

import (
	"database/sql"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dbutil"
	_ "github.com/mattn/go-sqlite3"
)

// goal: build table of mdl_field,default value
// uses: eph_named_default(kind, field, value): for the user's requested defaults.
//       mdl_kind(kind, path): for hierarchy.
//       mdl_field(kind, field, type)
// considerations:
// . property's actual kind
// . contradiction in specified values
// . contradiction in specified value vs type type ( alt: implicit conversion )
// . missing properties ( kind, field pair doesn't exist in model )
// o certainties: usually, seldom, never, always.
// o misspellings, near spellings ( ex. for missing fields )
func DetermineDefaults(m *Modeler, db *sql.DB) (err error) {
	// grab out the *actual* kind,field pairings
	var fieldType string
	var curr defaultInfo
	var list []defaultInfo
	// collect mdl_field rowid, type, and value.
	if e := dbutil.QueryAll(db, `
		/* table of hierarchy */
		with tree(kind,path,field,value) as 
		(select first.kind, first.path, p.field, p.value  
		    /* seed the search for kinds with the requested ephemera */
		    from eph_named_default p left join mdl_kind first
			on first.kind=p.kind 
			union all
			select super.kind, super.path, field, value 
			from tree kid, mdl_kind super
			where super.kind = substr(kid.path,0,instr(kid.path||",", ",")) 
		)
		select mf.rowid, mf.type, tree.value 
			from tree join mdl_field mf 
		   /* we are filtering kinds where the field in question is declared.
		      it'd be nicer to do this when we seed */
			where mf.field= tree.field 
			and mf.kind = tree.kind;`,
		func() (err error) {
			list = append(list, curr)
			return
		},
		&curr.idModelField, &fieldType, &curr.value,
	); e != nil {
		err = e
	} else {
		for _, n := range list {
			if e := m.WriteDefault(n.idModelField, n.value); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	return
}
