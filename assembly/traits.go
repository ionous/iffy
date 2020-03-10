package assembly

import (
	"database/sql"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/tables"
)

// goal: build table of trait, aspect, rank.
// considerations:
// o ambiguous ranks ( ignoring ranks for now. )
// . conflicting traits ( different aspects; future, resolve via implications? )
// . missing traits ( named but not specified )
// . missing aspects ( named but not specified )
// o misspellings, near spellings ( ex. for missing traits )
func DetermineAspects(m *Modeler, db *sql.DB) (err error) {
	var curr, last aspectInfo
	var traits []aspectInfo // cant read and write to the db simultaneously
	if e := tables.QueryAll(db, `select nt.name, na.name
	from eph_trait t join eph_named nt
		on (t.idNamedTrait = nt.rowid)
	left join eph_named na
		on (t.idNamedAspect = na.rowid)
	order by nt.name, na.name`, func() (err error) {
		switch traitsMatch, aspectsMatch := last.Trait == curr.Trait, last.Aspect == curr.Aspect; {
		case traitsMatch && !aspectsMatch:
			err = errutil.New("same trait different aspect", curr.Trait, curr.Aspect, last.Aspect)

		case !traitsMatch:
			traits = append(traits, curr)
			last = curr
		}
		return
	}, &curr.Trait, &curr.Aspect); e != nil {
		err = e
	} else {
		for _, t := range traits {
			// rank is not set yet.
			if e := m.WriteTrait(t.Aspect, t.Trait, 0); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	return
}

type aspectInfo struct {
	Aspect string
	Trait  string
}
